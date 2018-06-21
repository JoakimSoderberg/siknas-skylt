// Some real-time FFT! This visualizes music in the frequency domain using a
// polar-coordinate particle system. Particle size and radial distance are modulated
// using a filtered FFT. Color is sampled from an image.

import ddf.minim.analysis.*;
import ddf.minim.*;

OPC opc;
String host = "127.0.0.1";
int port = 7890;

PImage dot;
PImage colors;
Minim minim;
//AudioPlayer sound;
AudioInput sound;
FFT fft;
float[] fftFilter;

String filename = "data/groove.mp3";

float spin = 0.001;
float radiansPerBucket = radians(2);
float decay = 0.97;
float opacity = 50;
float minSize = 0.1;
float sizeScale = 0.6;

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

void load_layout(int x_offset, int y_offset, float scale)
{
  JSONArray points = loadJSONArray("layout.json");
  
  for (int j = 0; j < points.size(); j++) {
    JSONObject o = points.getJSONObject(j);
    JSONArray p = o.getJSONArray("point");
    opc.led(j, x_offset + (int)(p.getFloat(0)*scale), y_offset + (int)(p.getFloat(1)*scale));
  }
}

void setup()
{
  size(200, 200, P3D);

  minim = new Minim(this);

  // Small buffer size!
  //sound = minim.loadFile(filename, 512);
  //sound.loop();

  sound = minim.getLineIn(Minim.STEREO, 512);

  fft = new FFT(sound.bufferSize(), sound.sampleRate());
  fftFilter = new float[fft.specSize()];

  dot = loadImage("dot.png");
  colors = loadImage("colors.png");

  connect_opc();

  load_layout(0, 0, 200);
}

void draw()
{
  background(0);

  fft.forward(sound.mix);
  for (int i = 0; i < fftFilter.length; i++) {
    fftFilter[i] = max(fftFilter[i] * decay, log(1 + fft.getBand(i)));
  }
  
  for (int i = 0; i < fftFilter.length; i += 3) {   
    color rgb = colors.get(int(map(i, 0, fftFilter.length-1, 0, colors.width-1)), colors.height/2);
    tint(rgb, fftFilter[i] * opacity);
    blendMode(ADD);
 
    float size = height * (minSize + sizeScale * fftFilter[i]);
    PVector center = new PVector(width * (fftFilter[i] * 0.2), 0);
    center.rotate(millis() * spin + i * radiansPerBucket);
    center.add(new PVector(width * 0.5, height * 0.5));
 
    image(dot, center.x - size/2, center.y - size/2, size, size);
  }
}
