﻿using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;
using System.Text.RegularExpressions;
using Newtonsoft.Json;

namespace image_gui
{
    public partial class Form : System.Windows.Forms.Form
    {
        List<Point> ledLocations = new List<Point>();
        String layoutJson = "[]";
        FCServerConfig config = new FCServerConfig();

        public Form()
        {
            InitializeComponent();
            this.AddFadeCandy();
        }

        private void buttonOpen_Click(object sender, EventArgs e)
        {
            var openFile = new OpenFileDialog();
            openFile.Filter = "JPEG|*.jpg|PNG|*.png|GIF|*.gif";
            if (openFile.ShowDialog() == DialogResult.OK)
            {
                this.imageBox.Image = Image.FromFile(openFile.FileName);
                this.imageBox.ZoomToFit();
            }
        }
 
        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            Rectangle shape;
            var g = e.Graphics;
            int w = Math.Min((int)Math.Ceiling(20 / imageBox.ZoomFactor), 20);
            float penSize = Math.Min((int)Math.Ceiling(5 / imageBox.ZoomFactor), 3);

            // draw an outline around the top
            shape = imageBox.GetOffsetRectangle(50, 20, 100, 20);
            using (Pen pen = new Pen(Color.DarkRed, penSize))
            {
                DrawChannel(g, w, pen, this.config.Leds);
            }
        }

        private void DrawChannel(Graphics g, int w, Pen pen, List<Point> leds)
        {
            int i = 0;

            // Draw lines between the LEDs.
            if (leds.Count >= 2)
            {
                var points = new Point[leds.Count];
                i = 0;
                foreach (var led in leds)
                {
                    points[i] = imageBox.GetOffsetPoint(led);
                    i++;
                }

                g.DrawLines(pen, points);
            }

            // Draw the LEDs.
            i = 1;
            foreach (var led in leds)
            {
                var rect = new Rectangle(led.X - w / 2, led.Y - w / 2, w, w);
                rect = imageBox.GetOffsetRectangle(rect);
                g.FillEllipse(Brushes.Azure, rect);
                g.DrawEllipse(Pens.RoyalBlue, rect);
                g.DrawString(i.ToString(), SystemFonts.StatusFont, Brushes.Black, rect);
                i++;
            }
        }

        private void buttonClear_Click(object sender, EventArgs e)
        {
            var result = MessageBox.Show("Are you sure you want to clear all points?", "Work will be lost!",
                                        MessageBoxButtons.YesNo,
                                        MessageBoxIcon.Warning,
                                        MessageBoxDefaultButton.Button2);
            if (result == DialogResult.Yes)
            {
                this.config.Clear(false);
                this.imageBox.Refresh();
                this.GenerateJson();
            }
        }

        private void AddFadeCandy()
        {
            this.config.AddFadeCandy("fc" + this.treeView.Nodes.Count);
            this.treeView.Nodes.Clear();

            foreach (FadeCandy fc in this.config.devices)
            {
                var node = new TreeNode(fc.ToString());
                node.Tag = fc;
                this.treeView.Nodes.Add(node);

                foreach (FadeCandyChannel ch in fc.ledChannels)
                {
                    var childNode = new TreeNode(ch.ToString());
                    childNode.Tag = ch;
                    node.Nodes.Add(childNode);
                }
            }
        }

        private void AddPoint(Point loc)
        {
            this.ledLocations.Add(loc);
            this.config.AddLed(loc);

            this.imageBox.Refresh();
            this.GenerateJson();
        }

        private void deleteLastPoint()
        {
            this.config.RemoveLastLed();
            this.imageBox.Refresh();
            this.GenerateJson();
        }

        /*
         config.mapPixel = function (device, index) {
            //
            // Append a single device pixel to the mapping, returning the new OPC
            // pixel index. Consolidates contiguous mappings.
            //
            // Only supports channel 0 mappings of the form:
            // [ OPC Channel, First OPC Pixel, First output pixel, Pixel count ]
            //

        var devMap = config.mapDevice(device).map;
        var opcIndex = config.opcPixelCount++;
        var last = devMap[devMap.length - 1];

            if (last && last.length == 4
                && last[1] + last[3] == opcIndex
                && last[2] + last[3] == index) {
                // We can extend the last mapping
                last[3]++;
            } else {
                // New mapping line
                devMap.push([0, opcIndex, index, 1]);
            }

            return opcIndex;
        }
        */

        public class FadeCandyChannel
        {
            string name;
            public List<Point> leds = new List<Point>();
            FadeCandy parent;

            public FadeCandyChannel(string name, FadeCandy parent)
            {
                this.name = name;
                this.parent = parent;
            }

            public override string ToString()
            {
                return this.name;
            }
        }

        public class FadeCandy
        {
            string serial;
            const string type = "fadecandy";
            public List<FadeCandyChannel> ledChannels = new List<FadeCandyChannel>();
            private int selectedChannel;
            public int SelectedChannel { get => selectedChannel; set => selectedChannel = value; }

            public FadeCandy(string serial)
            {
                this.serial = serial;
                this.SelectedChannel = 0;
                this.ledChannels.Add(new FadeCandyChannel("1", this));
            }

            public List<Point> Leds
            {
                get
                {
                    var leds = new List<Point>();
                    foreach (FadeCandyChannel ch in this.ledChannels)
                    {
                        leds.AddRange(ch.leds);
                    }

                    return leds;
                }
            }

            public override string ToString()
            {
                return this.serial;
            }
        }

        public class FCServerConfig
        {
            bool verbose = true;

            public List<FadeCandy> devices = new List<FadeCandy>();
            int selectedIndex = 0;

            public FCServerConfig()
            {
            }

            public FadeCandy AddFadeCandy(string serial)
            {
                var fc = new FadeCandy("fc0");
                this.devices.Add(fc);
                return fc;
            }

            private List<Point> getCurrentLeds()
            {
                var fc = this.devices[this.selectedIndex];
                return fc.ledChannels[fc.SelectedChannel].leds;
            }

            public void Clear(bool all)
            {
                if (all)
                {
                    foreach (FadeCandy fc in this.devices)
                    {
                        foreach (FadeCandyChannel ch in fc.ledChannels)
                        {
                            ch.leds.Clear();
                        }
                    }
                }
                else
                {
                    this.getCurrentLeds().Clear();
                }
            }

            public void AddLed(Point p)
            {
                var leds = getCurrentLeds();
                leds.Add(p);
            }

            public void RemoveLastLed()
            {
                var leds = getCurrentLeds();
                if (leds.Count > 0)
                    leds.Remove(leds.Last());
            }

            public List<Point> Leds
            {
                get
                {
                    var leds = new List<Point>();

                    foreach (FadeCandy fc in this.devices)
                    {
                        leds.AddRange(fc.Leds);
                    }
 
                    return leds;
                }
            }
        }

        public class LayoutConfigItem
        {
            public int[] point;

            public LayoutConfigItem(Point p)
            {
                this.point = new int[3];
                this.point[0] = p.X;
                this.point[1] = p.Y;
                this.point[2] = 0;
            }
        }

        public class LayoutConfig
        {
            public List<LayoutConfigItem> points = new List<LayoutConfigItem>();
        }

        private void GenerateJson()
        {
            var fcserver = new FCServerConfig();
            LayoutConfig layout = new LayoutConfig
            {
                points = new List<LayoutConfigItem>(this.config.Leds.Count)
            };


            // TODO: We want one section 
            int i = 0;
            foreach (var p in this.config.Leds)
            {
                layout.points.Add(new LayoutConfigItem(p));
            }

            // The layout.json is just an array in the root of the JSON document.
            this.layoutJson = JsonConvert.SerializeObject(layout, Formatting.Indented);
            this.layoutJson = this.layoutJson.TrimEnd("}".ToCharArray());
            this.layoutJson = Regex.Replace(this.layoutJson, "\\s*{\\s*\\\"points\\\":\\s*", "", RegexOptions.Multiline);
            this.textBoxLayout.Text = this.layoutJson;
        }

        private void buttonRemove_Click(object sender, EventArgs e)
        {
            this.deleteLastPoint();
        }

        private void imageBox_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.KeyData == Keys.Delete)
            {
                this.deleteLastPoint();
            }
        }

        bool isPanning = false;

        private void imageBox_MouseClick(object sender, MouseEventArgs e)
        {
            if (isPanning)
                return;

            if (e.Button == MouseButtons.Left)
            {
                var loc = this.imageBox.PointToImage(e.Location);
                AddPoint(loc);
            }
            else
            {
                this.deleteLastPoint();
            }
        }
        
        private void imageBox_PanStart(object sender, EventArgs e)
        {
            this.isPanning = true;
        }

        private void imageBox_PanEnd(object sender, EventArgs e)
        {
            this.isPanning = false;
        }

        private void SaveFile(string text, string filename)
        {
            var dialog = new SaveFileDialog();
            dialog.Filter = "JSON|*.json|All files|*.txt";
            dialog.FileName = filename;

            if (dialog.ShowDialog() == DialogResult.OK)
            {
                var stream = dialog.OpenFile();
                if (stream != null)
                {
                    var msg = Encoding.UTF8.GetBytes(text);
                    stream.Write(msg, 0, msg.Length);
                    stream.Close();
                }
            }
        }

        private void buttonSaveLayoutAs_Click(object sender, EventArgs e)
        {
            this.GenerateJson();
            this.SaveFile(this.layoutJson, "layout.json");
        }
    }
}
