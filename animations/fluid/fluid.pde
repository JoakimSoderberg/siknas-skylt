// Some real-time FFT! This visualizes music in the frequency domain using a
// polar-coordinate particle system. Particle size and radial distance are modulated
// using a filtered FFT. Color is sampled from an image.

import ddf.minim.analysis.*;
import ddf.minim.*;

import ComputationalGeometry.*;

// FLUID SIMULATION EXAMPLE
import com.thomasdiewald.pixelflow.java.DwPixelFlow;
import com.thomasdiewald.pixelflow.java.fluid.DwFluid2D;

import javax.sound.sampled.AudioSystem;
import javax.sound.sampled.Mixer;


OPC opc;
String host = "192.168.1.71";
//String host = "192.168.99.100";
int port = 7899;

Minim minim;
AudioInput sound;
FFT fft;
float[] fftFilter;
BeatDetect beat;

PShape siknasMask;

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

// fluid simulation
DwFluid2D fluid;

// render target
PGraphics2D pg_fluid;

PImage colors;
float decay = 0.97;
float opacity = 50;

int prevTime;

PVector[] pts;
PVector prevPoint;
float burstDelay;
int colorIndex;

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
  
  siknasMask = loadShape("siknas-mask-cut-extra.svg");
  colors = loadImage("colors.png");
  
  // library context
  DwPixelFlow context = new DwPixelFlow(this);

  // fluid simulation
  fluid = new DwFluid2D(context, width/2, height/2, 1);

  // some fluid parameters
  fluid.param.dissipation_velocity = 0.70f;
  fluid.param.dissipation_density  = 0.99f;
  fluid.param.dissipation_temperature = 0.70f;
  fluid.param.vorticity               = 0.50f;
 
  // Used to determine how long between the fluid bursts.
  prevTime = millis();
  prevPoint = pts[0];
  burstDelay = 50;
  colorIndex = 0;

  // adding data to the fluid simulation
  fluid.addCallback_FluiData(new  DwFluid2D.FluidData() {
    public void update(DwFluid2D fluid) {

      beat.detect(sound.mix);
      fft.forward(sound.mix);

      for (int i = 0; i < fftFilter.length; i++) {
        fftFilter[i] = max(fftFilter[i] * decay, log(1 + fft.getBand(i)));
      }
      colorIndex = colorIndex % fftFilter.length;
      color rgb = colors.get(int(map(colorIndex, 0, fftFilter.length-1, 0, colors.width-1)), colors.height/2);
      colorIndex++;
      //tint(rgb, fftFilter[colorIndex] * opacity);*/
      //blendMode(ADD);

      // Allow manual interaction.
      if (mousePressed) {
        float px     = mouseX;
        float py     = height-mouseY;
        float vx     = (mouseX - pmouseX) * +15;
        float vy     = (mouseY - pmouseY) * -15;
        fluid.addVelocity(px, py, 14, vx, vy);
        fluid.addDensity (px, py, 20, 0.0f, 0.4f, 1.0f, 1.0f);
        fluid.addDensity (px, py, 8, 1.0f, 1.0f, 1.0f, 1.0f);
        fluid.addTemperature(px, py, 10.0f, 20.0f);
      }

      // Generate bursts periodically
      //if (millis() - prevTime >= burstDelay) {
      float intensity = noise(millis()) * 0.5;
      if (beat.isOnset()) {
        intensity = 1.0f;
        print(".");
      }

      /*if (millis() - prevTime >= burstDelay)*/ {
        PVector point = pts[int(random(pts.length))];
        float px = point.x;
        float py = height - point.y;
        float vx = (point.x - prevPoint.x) * 2;
        float vy = (point.y - prevPoint.y) * -2;
        fluid.addVelocity(px, py, intensity * 14, vx, vy);

        // Use the randomized color.
        float r = red(rgb) / 255.0f;
        float g = green(rgb) / 255.0f;
        float b = blue(rgb) / 255.0f;
        fluid.addDensity (px, py, 20, 0.0f, g, b, intensity);
        fluid.addDensity (px, py, 8, r, g, b, intensity);
        //fluid.addDensity (px, py, 20, 0.0f, 0.4f, 1.0f, 1.0f);
        //fluid.addDensity (px, py, 8, 1.0f, 1.0f, 1.0f, 1.0f);
        //fluid.addTemperature(px, py, 10.0f, 1.0f);

        prevPoint = point;
        prevTime = millis();
      }
    }
  });

  pg_fluid = (PGraphics2D) createGraphics(width, height, P2D);
}

float t = 0;

void draw()
{ 
  //shape(siknasMask, 16, 4, 175, 184);
  
  // update simulation
  fluid.update();

  // clear render target
  pg_fluid.beginDraw();
  pg_fluid.background(0);
  pg_fluid.endDraw();

  // render fluid stuff
  fluid.renderFluidTextures(pg_fluid, 0);

  // display
  image(pg_fluid, 0, 0);
}
