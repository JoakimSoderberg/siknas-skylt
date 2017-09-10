OPC opc;
PImage dot;

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
  size(800, 400);

  dot = loadImage("dot.png");

  // Connect to the local instance of fcserver. You can change this line to connect to another computer's fcserver
  opc = new OPC(this, "127.0.0.1", 7890);

  load_layout(0, 0, 400);
}

void draw()
{
  background(0);
  
  // Change the dot size as a function of time, to make it "throb"
  float dotSize = height * 0.6 * (1.0 + 0.2 * sin(millis() * 0.01));
  
  // Draw it centered at the mouse location
  image(dot, mouseX - dotSize/2, mouseY - dotSize/2, dotSize, dotSize);
}