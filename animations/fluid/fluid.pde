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
/*AudioInput sound;
FFT fft;
float[] fftFilter;*/

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

int prevTime;

PVector[] pts;
PVector prevPoint;

void setup()
{
  size(200, 200, P2D);

  /*
  // TODO: Enable setting mixer via command line arguments
  minim = new Minim(this);
  
  Mixer.Info[] mixerInfo;
  mixerInfo = AudioSystem.getMixerInfo(); 
  
  for(int i = 0; i < mixerInfo.length; i++) {
    print(i + ": " + mixerInfo[i].getName() + "\n");
  } 
  // 0 is pulseaudio mixer on GNU/Linux
  Mixer mixer = AudioSystem.getMixer(mixerInfo[2]);
  minim.setInputMixer(mixer);

  sound = minim.getLineIn(Minim.STEREO, 512);

  fft = new FFT(sound.bufferSize(), sound.sampleRate());
  fftFilter = new float[fft.specSize()];
  */

  connect_opc();

  pts = load_layout(0, 0, 200);
  
  siknasMask = loadShape("siknas-mask-cut-extra.svg");
  
  // library context
  DwPixelFlow context = new DwPixelFlow(this);

  // fluid simulation
  fluid = new DwFluid2D(context, width/2, height/2, 1);

  // some fluid parameters
  fluid.param.dissipation_velocity = 0.70f;
  fluid.param.dissipation_density  = 0.99f;
 
  prevTime = millis();
  prevPoint = pts[0];

  // adding data to the fluid simulation
  fluid.addCallback_FluiData(new  DwFluid2D.FluidData() {
    public void update(DwFluid2D fluid) {
      if (mousePressed) {
        float px     = mouseX;
        float py     = height-mouseY;
        float vx     = (mouseX - pmouseX) * +15;
        float vy     = (mouseY - pmouseY) * -15;
        fluid.addVelocity(px, py, 14, vx, vy);
        fluid.addDensity (px, py, 20, 0.0f, 0.4f, 1.0f, 1.0f);
        fluid.addDensity (px, py, 8, 1.0f, 1.0f, 1.0f, 1.0f);
      }

      if (millis() - prevTime >= 1000) {
        PVector point = pts[int(random(pts.length))];
        float px = point.x;
        float py = height - point.y;
        float vx = (point.x - prevPoint.x) * 15;
        float vy = (point.y - prevPoint.y) * -15;
        fluid.addVelocity(px, py, 14, vx, vy);
        fluid.addDensity (px, py, 20, 0.0f, 0.4f, 1.0f, 1.0f);
        fluid.addDensity (px, py, 8, 1.0f, 1.0f, 1.0f, 1.0f);

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
