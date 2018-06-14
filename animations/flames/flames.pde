OPC opc;
PImage im;
String host = "127.0.0.1";
int port = 7890;

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

  // Load a sample image
  im = loadImage("flames.jpeg");

  if (args != null) {
    String[] parts = split(args[0], ":");
    host = parts[0];
    port = Integer.parseInt(parts[1]);
  }

  println("Connecting to ", host, ":", port);

  // Connect to the local instance of fcserver
  opc = new OPC(this, host, port);

  load_layout(0,0, 200);
}

int imHeight;
float speed;
float y;

void draw()
{
  // Scale the image so that it matches the width of the window
  imHeight = im.height * width / im.width;

  // Scroll down slowly, and wrap around
  speed = 0.05;
  y = (millis() * -speed) % imHeight;
  
  // Use two copies of the image, so it seems to repeat infinitely  
  image(im, 0, y, width, imHeight);
  image(im, 0, y + imHeight, width, imHeight);
}
