OPC opc;
PImage dot1, dot2;

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
  size(200, 200);
  frameRate(20);

  dot1 = loadImage("greenDot.png");
  dot2 = loadImage("purpleDot.png");

  // Connect to the local instance of fcserver. You can change this line to connect to another computer's fcserver
  opc = new OPC(this, "127.0.0.1", 7890);

  load_layout(0, 0, 200);
}

float px, py;

void draw()
{
  background(0);
  blendMode(ADD);
  
  // Smooth out the mouse location
  px += (mouseX - px) * 0.1;
  py += (mouseY - py) * 0.1;

  float a = millis() * 0.001;
  float r = py * 0.5;
  float dotSize = r * 4;  

  float dx = width/2 + cos(a) * r;
  float dy = height/2 + sin(a) * r;
  
  // Draw it centered at the mouse location
  image(dot1, dx - dotSize/2, dy - dotSize/2, dotSize, dotSize);

  // Another dot, mirrored around the center
  image(dot2, width - dx - dotSize/2, height - dy - dotSize/2, dotSize, dotSize);
}