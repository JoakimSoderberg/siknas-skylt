// Particle system with multiple attraction points.
// Spawns centered around a random point, lives out a cycle and dies; the cycle repeats.

int numParticles = 10;
float cornerCoefficient = 0.2;
int integrationSteps = 20;
float maxOpacity = 100;
float stepFast = 1.0 / 100;
float stepSlow = 1.0 / 3000;
float energyThreshold = 10.0;
float brightnessThreshold = 0.8;

OPC opc;
PImage dot;
PImage colors;
Particle[] particles;
PVector[] corners;
float epoch = 0;

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
  size(600, 600, P3D);
  frameRate(30);

  dot = loadImage("dot.png");
  colors = loadImage("colors.png");
  colors.loadPixels();

  // Connect to the local instance of fcserver
  opc = new OPC(this, "127.0.0.1", 7890);

  load_layout(0, 0, 600);

  // Attraction points
  corners = new PVector[3];
  corners[0] = new PVector(width * 1/4, height * 0.5);
  corners[1] = new PVector(width * 2/4, height * 0.5);
  corners[2] = new PVector(width * 3/4, height * 0.5);
 
  beginEpoch();
}

void beginEpoch()
{
  epoch = 0;
 
  // Center of bundle
  float s = 0.5;
  float cx = width * (0.5 + random(-s, s));
  float cy = height * (0.5 + random(-s, s));
 
  // Half-width of particle bundle
  float w = width * 0.2;
 
  particles = new Particle[numParticles];
  for (int i = 0; i < particles.length; i++) {
    color rgb = colors.pixels[int(random(0, colors.width * colors.height))];
    particles[i] = new Particle(
      cx + random(-w, w),
      cy + random(-w, w), rgb);
  }
}

void draw()
{
  background(0);
  
  // How much energy is still left?
  float energy = 0;
  for (int i = 0; i < particles.length; i++) {
    energy += particles[i].energy();
  }
  
  // How bright is our brightest pixel?
  float brightness = 0;
  for (int i = 0; i < opc.pixelLocations.length; i++) {
    color rgb = opc.getPixel(i);
    brightness = max(brightness, max(red(rgb), max(blue(rgb), green(rgb))));
  }
  brightness /= 255.0;
  
  text("Energy: " + energy, 2, 12);
  text("Brightness: " + brightness, 2, 25);

  // What's interesting? Can we maintain high brightness and high energy?
  // These are normally conflicting goals. If we've managed to balance the two,
  // keep going to see how it turns out.
  if (energy > energyThreshold && brightness > brightnessThreshold) {
 
    // Time moves slower when we're interested
    epoch += stepSlow;
    text("+", 2, 40);
  } else {
    epoch += stepFast;
  }
  
  if (epoch > 1) {
    beginEpoch();
  }
    
  for (int step = 0; step < integrationSteps; step++) {
    for (int i = 0; i < particles.length; i++) {
      particles[i].integrate();

      // Each particle is attracted by the corners
      for (int j = 0; j < corners.length; j++) {
        particles[i].attract(corners[j], cornerCoefficient);
      }
    }
  }

  for (int i = 0; i < particles.length; i++) {
    particles[i].draw(sin(epoch * PI) * maxOpacity);
  }
}