// Some real-time FFT! This visualizes music in the frequency domain using a
// polar-coordinate particle system. Particle size and radial distance are modulated
// using a filtered FFT. Color is sampled from an image.

import ddf.minim.analysis.*;
import ddf.minim.*;

import ComputationalGeometry.*;

import javax.sound.sampled.AudioSystem;
import javax.sound.sampled.Mixer;


OPC opc;
//String host = "192.168.1.71";
String host = "192.168.99.100";
//int port = 7899;
int port = 7890;

Minim minim;
AudioInput sound;
FFT fft;
float[] fftFilter;
BeatDetect beat;

PShape siknasMask;
PShape[] siknas = new PShape[6];

void connect_opc()
{
  if (args != null) {
    String[] parts = split(args[0], ":");
    host = parts[0];
    port = Integer.parseInt(parts[1]);
  }

  println("Connecting to ", host, ":", port);

  opc = new OPC(this, host, port);
}

PVector[] load_layout(int x_offset, int y_offset, float scale)
{
  JSONArray points = loadJSONArray("data/layout.json");
  PVector[] scaled_points = new PVector[points.size()];
  
  for (int j = 0; j < points.size(); j++) {
    JSONObject o = points.getJSONObject(j);
    JSONArray p = o.getJSONArray("point");
    int x = x_offset + (int)(p.getFloat(0)*scale);
    int y = y_offset + (int)(p.getFloat(1)*scale);
    opc.led(j, x, y);

    scaled_points[j] = new PVector(x, y);
  }

  return scaled_points;
}

PImage colors;

PImage ripplesImg;

PVector[] pts;

int letterIndex = 0;
color[] siknasColors = new color[6];
float[] siknasFade = new float[6];
float xoff = 0.0;

// Backround ripples.
int imgHeight;
float imgSpeed = 0.01;
float y;

// Fade speed on beats.
float fadeSpeed = 0.025;

void setup()
{
  size(200, 200, P2D);

  minim = new Minim(this);
  
  // TODO: Enable setting mixer via command line arguments.
  // TODO: Make this work properly on Raspberry Pi.
  /*Mixer.Info[] mixerInfo;
  mixerInfo = AudioSystem.getMixerInfo(); 
  
  for(int i = 0; i < mixerInfo.length; i++) {
    print(i + ": " + mixerInfo[i].getName() + "\n");
  } 
  // 0 is pulseaudio mixer on GNU/Linux
  Mixer mixer = AudioSystem.getMixer(mixerInfo[2]);
  minim.setInputMixer(mixer);*/

  sound = minim.getLineIn(Minim.STEREO, 512);

  fft = new FFT(sound.bufferSize(), sound.sampleRate());
  fftFilter = new float[fft.specSize()];
  
  beat = new BeatDetect(sound.bufferSize(), sound.sampleRate());
  beat.setSensitivity(50);
  beat.detectMode(BeatDetect.SOUND_ENERGY);

  connect_opc();

  pts = load_layout(0, 0, 200);
  
  //siknasMask = loadShape("siknas-mask-cut-extra.svg");
  
  siknas[0] = loadShape("siknas-s-only.svg");
  siknas[1] = loadShape("siknas-i-only.svg");
  siknas[2] = loadShape("siknas-k-only.svg");
  siknas[3] = loadShape("siknas-n-only.svg");
  siknas[4] = loadShape("siknas-a-only.svg");
  siknas[5] = loadShape("siknas-sb-only.svg");

  colors = loadImage("colors.png");
  ripplesImg = loadImage("ripples.png");
}

void draw()
{
  background(0);

  imgHeight = ripplesImg.height * width / ripplesImg.width;
  y = (millis() * imgSpeed) % imgHeight;

  // Use two copies of the image, so it seems to repeat infinitely  
  image(ripplesImg, 0, y, width, imgHeight);
  image(ripplesImg, 0, y - imgHeight, width, imgHeight);

  beat.detect(sound.mix);
  xoff += 0.1;

  // On a beat we want to light up a character.
  if (beat.isOnset()) {
    letterIndex = (int)(round(noise(xoff) * 16)) - 1;
    print(letterIndex);
    letterIndex %= 6;
  }

  for (int i = 0; i < siknas.length; i++) {
    siknas[i].disableStyle();
  }

  // Hand tweaked magic to position letters over LEDs.
  translate(29, 18);
  scale(0.1345, 0.19);

  for (int i = 0; i < siknas.length; i++) {
    // Align the last three letters a bit better.
    if (i == 3)
      translate(-10.0, 10.0);

    // The letter to light up on a beat.
    if (i == letterIndex) {
      if (beat.isOnset()) {
        // New color only on the detected beat, so we don't flash madly.
        //siknasColors[i] = color(random(255), random(255), random(255), 255);
        siknasFade[i] = 1.0;
        siknasColors[i] = colors.get((int)random(colors.width - 1), colors.height/2);
      }
    }

    // Lerp color towards black over time.
    siknasColors[i] = lerpColor(siknasColors[i], color(0, 0, 0, 0), 1.0 - siknasFade[i]);

    fill(siknasColors[i]);
    stroke(siknasColors[i]);
    
    // Avoid drawing black shapes.
    if (siknasFade[i] > 0.0 && brightness(siknasColors[i]) > 0.0) {
      shape(siknas[i]);
    }

    siknasFade[i] -= 0.005;
  }
}
