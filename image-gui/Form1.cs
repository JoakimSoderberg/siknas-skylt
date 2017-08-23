using System;
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

        public Form()
        {
            InitializeComponent();
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
                int i = 0;
                if (this.ledLocations.Count >= 2)
                {
                    var points = new Point[this.ledLocations.Count];
                    i = 0;
                    foreach (var led in this.ledLocations)
                    {
                        points[i] = imageBox.GetOffsetPoint(led);
                        i++;
                    }

                    g.DrawLines(pen, points);
                }

                i = 1;
                foreach (var led in this.ledLocations)
                {
                    var rect = new Rectangle(led.X - w / 2, led.Y - w / 2, w, w);
                    rect = imageBox.GetOffsetRectangle(rect);
                    g.FillEllipse(Brushes.Azure, rect);
                    g.DrawEllipse(Pens.RoyalBlue, rect);
                    g.DrawString(i.ToString(), SystemFonts.StatusFont, Brushes.Black, rect);
                    i++;
                }
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
                this.ledLocations.Clear();
                this.imageBox.Refresh();
                this.GenerateJson();
            }
        }

        private void AddPoint(Point loc)
        {
            this.ledLocations.Add(loc);
            this.imageBox.Refresh();
            this.GenerateJson();
        }

        private void deleteLastPoint()
        {
            if (this.ledLocations.Count > 0)
            {
                this.ledLocations.Remove(this.ledLocations.Last());
                this.imageBox.Refresh();
                this.GenerateJson();
            }
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

        internal class FCServerConfig
        {

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

        internal struct Point3
        {
            int x;
            int y;
            int z;
        }

        private void GenerateJson()
        {
            var fcserver = new FCServerConfig();
            LayoutConfig layout = new LayoutConfig();
            layout.points = new List<LayoutConfigItem>(this.ledLocations.Count);


            // TODO: We want one section 
            int i = 0;
            foreach (var p in this.ledLocations)
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
