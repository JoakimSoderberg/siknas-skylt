using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace image_gui
{
    public partial class Form : System.Windows.Forms.Form
    {
        List<Point> ledLocations = new List<Point>();

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
            }
        }

        private void deleteLastPoint()
        {
            if (this.ledLocations.Count > 0)
            {
                this.ledLocations.Remove(this.ledLocations.Last());
                this.imageBox.Refresh();
            }
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
                this.ledLocations.Add(loc);
                this.imageBox.Refresh();
            }
            else
            {
                this.deleteLastPoint();
            }
        }

        private void imageBox_DoubleClick(object sender, MouseEventArgs e)
        {
        }

        private void imageBox_PanStart(object sender, EventArgs e)
        {
            this.isPanning = true;
        }

        private void imageBox_PanEnd(object sender, EventArgs e)
        {
            this.isPanning = false;
        }
    }
}
