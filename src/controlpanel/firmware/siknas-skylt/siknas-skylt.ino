/*
 ESP8266 Blink by Simon Peter
 Blink the blue LED on the ESP-01 module
 This example code is in the public domain
 
 The blue LED on the ESP-01 module is connected to GPIO1 
 (which is also the TXD pin; so we cannot use Serial.print() at the same time)
 
 Note that this sketch uses BUILTIN_LED to find the pin with the internal LED
*/
#define BUILTIN_LED  11

#define METER_PIN A9

#define POT_PIN_R A0
#define POT_PIN_G A1
#define POT_PIN_B A2

#define ROTARY_PIN_START 1
#define ROTARY_PIN_END 4

#define R_COLOR_IDX 0
#define G_COLOR_IDX 1
#define B_COLOR_IDX 2

#define COLOR_READ_BUF 20

int color_vals[COLOR_READ_BUF][4];
int color_totals[4];
int color_read_idx;

void setup()
{
  pinMode(BUILTIN_LED, OUTPUT);     // Initialize the BUILTIN_LED pin as an output

  // Rotary switch inputs.
  for (int pin = ROTARY_PIN_START; pin <= ROTARY_PIN_END; pin++)
  {
    pinMode(pin, INPUT_PULLUP);  
  }
  
  pinMode(METER_PIN, OUTPUT);

  pinMode(POT_PIN_R, INPUT);
  pinMode(POT_PIN_G, INPUT);
  pinMode(POT_PIN_B, INPUT);

  Serial.begin(9600);
}

void printPotPin(int color_idx, char c, int val)
{
  //Serial.print(c);
  Serial.print(val);
  Serial.write(" ");
  #if 0
  Serial.write(" total: ");
  Serial.print(color_totals[color_idx]);
  #endif
  //Serial.write("\n");
}

int read_smooth_color(int color_idx, int pin)
{
  color_totals[color_idx] -= color_vals[color_read_idx][color_idx];

  // Note! We read these multiple times to get rid of cross talk between
  // the different analog channels.
  for (int i = 0; i < 4; i++)
  {
    color_vals[color_read_idx][color_idx] = analogRead(pin);
  }

  color_totals[color_idx] += color_vals[color_read_idx][color_idx];

  delay(1);

  return map(color_totals[color_idx] / COLOR_READ_BUF, 0, 1023, 0, 255);
}

// the loop function runs over and over again forever
void loop()
{
  unsigned long now = millis();
  static unsigned long print_time = 0;

  int program_pin = 0;

  int r = read_smooth_color(R_COLOR_IDX, POT_PIN_R);
  int b = read_smooth_color(B_COLOR_IDX, POT_PIN_B);
  int g = read_smooth_color(G_COLOR_IDX, POT_PIN_G);

  color_read_idx++;
  if (color_read_idx >= COLOR_READ_BUF)
  {
    color_read_idx = 0;
  }

  // The program selector.
  for (int pin = ROTARY_PIN_START; pin <= ROTARY_PIN_END; pin++)
  {
    if (!digitalRead(pin)) {
      program_pin = pin;
    }
  }

  analogWrite(METER_PIN, min((int)(program_pin * 255/4), 255));

  if ((now - print_time) > 1000)
  {
    print_time = millis();
    printPotPin(0, 'P', program_pin);
    printPotPin(0, 'R', r);
    printPotPin(1, 'G', g);
    printPotPin(2, 'B', b);
    printPotPin(0, 'B', 255); // TODO: Write brightness pot
    //Serial.write("-------------\n");
    Serial.write("\n");
  }

  
}
