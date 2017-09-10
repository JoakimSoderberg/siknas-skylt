OPC opc;
PImage im;

void load_layout()
{
  JSONArray points = loadJSONArray("layout.json");
  
  for (int j = 0; j < points.size(); j++) {
    JSONObject o = points.getJSONObject(j);
    JSONArray p = o.getJSONArray("point");
    opc.led(j, (int)(p.getFloat(0)*200), (int)(p.getFloat(1)*200));
  }
}

void setup()
{
  size(200, 200);

  // Load a sample image
  im = loadImage("flames.jpeg");

  // Connect to the local instance of fcserver
  opc = new OPC(this, "192.168.1.86", 7890);

  // Map one 64-LED strip to the center of the window
  //opc.ledStrip(0, 64, width/2, height/2, width / 70.0, 0, false);

  load_layout();
}

void draw()
{
  // Scale the image so that it matches the width of the window
  int imHeight = im.height * width / im.width;

  // Scroll down slowly, and wrap around
  float speed = 0.05;
  float y = (millis() * -speed) % imHeight;
  
  // Use two copies of the image, so it seems to repeat infinitely  
  image(im, 0, y, width, imHeight);
  image(im, 0, y + imHeight, width, imHeight);
}