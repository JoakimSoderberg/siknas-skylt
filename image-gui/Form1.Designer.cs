namespace image_gui
{
    partial class Form
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.panel1 = new System.Windows.Forms.Panel();
            this.buttonSaveAs = new System.Windows.Forms.Button();
            this.buttonOpen = new System.Windows.Forms.Button();
            this.tabControl = new System.Windows.Forms.TabControl();
            this.tabPageImage = new System.Windows.Forms.TabPage();
            this.imageBox = new Cyotek.Windows.Forms.ImageBox();
            this.tabPageJSON = new System.Windows.Forms.TabPage();
            this.textBox1 = new System.Windows.Forms.TextBox();
            this.groupBoxFile = new System.Windows.Forms.GroupBox();
            this.groupBoxPoints = new System.Windows.Forms.GroupBox();
            this.buttonClear = new System.Windows.Forms.Button();
            this.buttonRemove = new System.Windows.Forms.Button();
            this.panel1.SuspendLayout();
            this.tabControl.SuspendLayout();
            this.tabPageImage.SuspendLayout();
            this.tabPageJSON.SuspendLayout();
            this.groupBoxFile.SuspendLayout();
            this.groupBoxPoints.SuspendLayout();
            this.SuspendLayout();
            // 
            // panel1
            // 
            this.panel1.Controls.Add(this.groupBoxPoints);
            this.panel1.Controls.Add(this.groupBoxFile);
            this.panel1.Dock = System.Windows.Forms.DockStyle.Bottom;
            this.panel1.Location = new System.Drawing.Point(0, 452);
            this.panel1.Name = "panel1";
            this.panel1.Size = new System.Drawing.Size(687, 112);
            this.panel1.TabIndex = 1;
            // 
            // buttonSaveAs
            // 
            this.buttonSaveAs.Location = new System.Drawing.Point(7, 51);
            this.buttonSaveAs.Name = "buttonSaveAs";
            this.buttonSaveAs.Size = new System.Drawing.Size(105, 23);
            this.buttonSaveAs.TabIndex = 1;
            this.buttonSaveAs.Text = "Save Json as...";
            this.buttonSaveAs.UseVisualStyleBackColor = true;
            // 
            // buttonOpen
            // 
            this.buttonOpen.Location = new System.Drawing.Point(7, 19);
            this.buttonOpen.Name = "buttonOpen";
            this.buttonOpen.Size = new System.Drawing.Size(105, 23);
            this.buttonOpen.TabIndex = 0;
            this.buttonOpen.Text = "Open Image...";
            this.buttonOpen.UseVisualStyleBackColor = true;
            this.buttonOpen.Click += new System.EventHandler(this.buttonOpen_Click);
            // 
            // tabControl
            // 
            this.tabControl.Controls.Add(this.tabPageImage);
            this.tabControl.Controls.Add(this.tabPageJSON);
            this.tabControl.Dock = System.Windows.Forms.DockStyle.Fill;
            this.tabControl.Location = new System.Drawing.Point(0, 0);
            this.tabControl.Name = "tabControl";
            this.tabControl.SelectedIndex = 0;
            this.tabControl.Size = new System.Drawing.Size(687, 452);
            this.tabControl.TabIndex = 2;
            // 
            // tabPageImage
            // 
            this.tabPageImage.Controls.Add(this.imageBox);
            this.tabPageImage.Location = new System.Drawing.Point(4, 22);
            this.tabPageImage.Name = "tabPageImage";
            this.tabPageImage.Padding = new System.Windows.Forms.Padding(3);
            this.tabPageImage.Size = new System.Drawing.Size(679, 426);
            this.tabPageImage.TabIndex = 0;
            this.tabPageImage.Text = "Image";
            this.tabPageImage.UseVisualStyleBackColor = true;
            // 
            // imageBox
            // 
            this.imageBox.AllowDoubleClick = true;
            this.imageBox.Dock = System.Windows.Forms.DockStyle.Fill;
            this.imageBox.Location = new System.Drawing.Point(3, 3);
            this.imageBox.Name = "imageBox";
            this.imageBox.ShowPixelGrid = true;
            this.imageBox.Size = new System.Drawing.Size(673, 420);
            this.imageBox.TabIndex = 0;
            this.imageBox.PanEnd += new System.EventHandler(this.imageBox_PanEnd);
            this.imageBox.PanStart += new System.EventHandler(this.imageBox_PanStart);
            this.imageBox.Paint += new System.Windows.Forms.PaintEventHandler(this.imageBox_Paint);
            this.imageBox.KeyDown += new System.Windows.Forms.KeyEventHandler(this.imageBox_KeyDown);
            this.imageBox.MouseClick += new System.Windows.Forms.MouseEventHandler(this.imageBox_MouseClick);
            this.imageBox.MouseDoubleClick += new System.Windows.Forms.MouseEventHandler(this.imageBox_DoubleClick);
            // 
            // tabPageJSON
            // 
            this.tabPageJSON.Controls.Add(this.textBox1);
            this.tabPageJSON.Location = new System.Drawing.Point(4, 22);
            this.tabPageJSON.Name = "tabPageJSON";
            this.tabPageJSON.Padding = new System.Windows.Forms.Padding(3);
            this.tabPageJSON.Size = new System.Drawing.Size(679, 487);
            this.tabPageJSON.TabIndex = 1;
            this.tabPageJSON.Text = "Json";
            this.tabPageJSON.UseVisualStyleBackColor = true;
            // 
            // textBox1
            // 
            this.textBox1.Dock = System.Windows.Forms.DockStyle.Fill;
            this.textBox1.Location = new System.Drawing.Point(3, 3);
            this.textBox1.Multiline = true;
            this.textBox1.Name = "textBox1";
            this.textBox1.ReadOnly = true;
            this.textBox1.Size = new System.Drawing.Size(673, 481);
            this.textBox1.TabIndex = 0;
            // 
            // groupBoxFile
            // 
            this.groupBoxFile.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.groupBoxFile.Controls.Add(this.buttonOpen);
            this.groupBoxFile.Controls.Add(this.buttonSaveAs);
            this.groupBoxFile.Location = new System.Drawing.Point(12, 11);
            this.groupBoxFile.Name = "groupBoxFile";
            this.groupBoxFile.Size = new System.Drawing.Size(131, 89);
            this.groupBoxFile.TabIndex = 2;
            this.groupBoxFile.TabStop = false;
            this.groupBoxFile.Text = "File";
            // 
            // groupBoxPoints
            // 
            this.groupBoxPoints.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.groupBoxPoints.Controls.Add(this.buttonRemove);
            this.groupBoxPoints.Controls.Add(this.buttonClear);
            this.groupBoxPoints.Location = new System.Drawing.Point(481, 11);
            this.groupBoxPoints.Name = "groupBoxPoints";
            this.groupBoxPoints.Size = new System.Drawing.Size(194, 89);
            this.groupBoxPoints.TabIndex = 3;
            this.groupBoxPoints.TabStop = false;
            this.groupBoxPoints.Text = "Points";
            // 
            // buttonClear
            // 
            this.buttonClear.Location = new System.Drawing.Point(6, 19);
            this.buttonClear.Name = "buttonClear";
            this.buttonClear.Size = new System.Drawing.Size(75, 23);
            this.buttonClear.TabIndex = 0;
            this.buttonClear.Text = "Clear all";
            this.buttonClear.UseVisualStyleBackColor = true;
            this.buttonClear.Click += new System.EventHandler(this.buttonClear_Click);
            // 
            // buttonRemove
            // 
            this.buttonRemove.Location = new System.Drawing.Point(6, 51);
            this.buttonRemove.Name = "buttonRemove";
            this.buttonRemove.Size = new System.Drawing.Size(75, 23);
            this.buttonRemove.TabIndex = 1;
            this.buttonRemove.Text = "Remove last";
            this.buttonRemove.UseVisualStyleBackColor = true;
            this.buttonRemove.Click += new System.EventHandler(this.buttonRemove_Click);
            // 
            // Form
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(687, 564);
            this.Controls.Add(this.tabControl);
            this.Controls.Add(this.panel1);
            this.MinimumSize = new System.Drawing.Size(300, 300);
            this.Name = "Form";
            this.Text = "Form";
            this.panel1.ResumeLayout(false);
            this.tabControl.ResumeLayout(false);
            this.tabPageImage.ResumeLayout(false);
            this.tabPageJSON.ResumeLayout(false);
            this.tabPageJSON.PerformLayout();
            this.groupBoxFile.ResumeLayout(false);
            this.groupBoxPoints.ResumeLayout(false);
            this.ResumeLayout(false);

        }

        #endregion
        private System.Windows.Forms.Panel panel1;
        private System.Windows.Forms.Button buttonOpen;
        private System.Windows.Forms.TabControl tabControl;
        private System.Windows.Forms.TabPage tabPageImage;
        private System.Windows.Forms.TabPage tabPageJSON;
        private System.Windows.Forms.TextBox textBox1;
        private System.Windows.Forms.Button buttonSaveAs;
        private Cyotek.Windows.Forms.ImageBox imageBox;
        private System.Windows.Forms.GroupBox groupBoxFile;
        private System.Windows.Forms.GroupBox groupBoxPoints;
        private System.Windows.Forms.Button buttonClear;
        private System.Windows.Forms.Button buttonRemove;
    }
}

