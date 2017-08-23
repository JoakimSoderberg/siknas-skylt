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
            this.groupBoxPoints = new System.Windows.Forms.GroupBox();
            this.buttonRemove = new System.Windows.Forms.Button();
            this.buttonClear = new System.Windows.Forms.Button();
            this.groupBoxFile = new System.Windows.Forms.GroupBox();
            this.buttonOpen = new System.Windows.Forms.Button();
            this.buttonSaveLayoutAs = new System.Windows.Forms.Button();
            this.tabControl = new System.Windows.Forms.TabControl();
            this.tabPageImage = new System.Windows.Forms.TabPage();
            this.imageBox = new Cyotek.Windows.Forms.ImageBox();
            this.tabPageLayout = new System.Windows.Forms.TabPage();
            this.textBoxLayout = new System.Windows.Forms.TextBox();
            this.treeView = new System.Windows.Forms.TreeView();
            this.groupBoxCurrent = new System.Windows.Forms.GroupBox();
            this.tabPageFcserver = new System.Windows.Forms.TabPage();
            this.textBoxFcserver = new System.Windows.Forms.TextBox();
            this.buttonSaveConfigAs = new System.Windows.Forms.Button();
            this.buttonAddFadecandy = new System.Windows.Forms.Button();
            this.buttonRemoveFadecandy = new System.Windows.Forms.Button();
            this.panel1.SuspendLayout();
            this.groupBoxPoints.SuspendLayout();
            this.groupBoxFile.SuspendLayout();
            this.tabControl.SuspendLayout();
            this.tabPageImage.SuspendLayout();
            this.tabPageLayout.SuspendLayout();
            this.tabPageFcserver.SuspendLayout();
            this.SuspendLayout();
            // 
            // panel1
            // 
            this.panel1.Controls.Add(this.buttonRemoveFadecandy);
            this.panel1.Controls.Add(this.buttonAddFadecandy);
            this.panel1.Controls.Add(this.groupBoxCurrent);
            this.panel1.Controls.Add(this.treeView);
            this.panel1.Controls.Add(this.groupBoxPoints);
            this.panel1.Controls.Add(this.groupBoxFile);
            this.panel1.Dock = System.Windows.Forms.DockStyle.Bottom;
            this.panel1.Location = new System.Drawing.Point(0, 429);
            this.panel1.Name = "panel1";
            this.panel1.Size = new System.Drawing.Size(687, 157);
            this.panel1.TabIndex = 1;
            // 
            // groupBoxPoints
            // 
            this.groupBoxPoints.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.groupBoxPoints.Controls.Add(this.buttonRemove);
            this.groupBoxPoints.Controls.Add(this.buttonClear);
            this.groupBoxPoints.Location = new System.Drawing.Point(588, 3);
            this.groupBoxPoints.Name = "groupBoxPoints";
            this.groupBoxPoints.Size = new System.Drawing.Size(87, 112);
            this.groupBoxPoints.TabIndex = 3;
            this.groupBoxPoints.TabStop = false;
            this.groupBoxPoints.Text = "Points";
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
            // groupBoxFile
            // 
            this.groupBoxFile.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.groupBoxFile.Controls.Add(this.buttonSaveConfigAs);
            this.groupBoxFile.Controls.Add(this.buttonOpen);
            this.groupBoxFile.Controls.Add(this.buttonSaveLayoutAs);
            this.groupBoxFile.Location = new System.Drawing.Point(12, 3);
            this.groupBoxFile.Name = "groupBoxFile";
            this.groupBoxFile.Size = new System.Drawing.Size(131, 112);
            this.groupBoxFile.TabIndex = 2;
            this.groupBoxFile.TabStop = false;
            this.groupBoxFile.Text = "File";
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
            // buttonSaveLayoutAs
            // 
            this.buttonSaveLayoutAs.Location = new System.Drawing.Point(7, 48);
            this.buttonSaveLayoutAs.Name = "buttonSaveLayoutAs";
            this.buttonSaveLayoutAs.Size = new System.Drawing.Size(105, 23);
            this.buttonSaveLayoutAs.TabIndex = 1;
            this.buttonSaveLayoutAs.Text = "Save layout as...";
            this.buttonSaveLayoutAs.UseVisualStyleBackColor = true;
            this.buttonSaveLayoutAs.Click += new System.EventHandler(this.buttonSaveLayoutAs_Click);
            // 
            // tabControl
            // 
            this.tabControl.Controls.Add(this.tabPageImage);
            this.tabControl.Controls.Add(this.tabPageLayout);
            this.tabControl.Controls.Add(this.tabPageFcserver);
            this.tabControl.Dock = System.Windows.Forms.DockStyle.Fill;
            this.tabControl.Location = new System.Drawing.Point(0, 0);
            this.tabControl.Name = "tabControl";
            this.tabControl.SelectedIndex = 0;
            this.tabControl.Size = new System.Drawing.Size(687, 429);
            this.tabControl.TabIndex = 2;
            // 
            // tabPageImage
            // 
            this.tabPageImage.Controls.Add(this.imageBox);
            this.tabPageImage.Location = new System.Drawing.Point(4, 22);
            this.tabPageImage.Name = "tabPageImage";
            this.tabPageImage.Padding = new System.Windows.Forms.Padding(3);
            this.tabPageImage.Size = new System.Drawing.Size(679, 403);
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
            this.imageBox.Size = new System.Drawing.Size(673, 397);
            this.imageBox.TabIndex = 0;
            this.imageBox.PanEnd += new System.EventHandler(this.imageBox_PanEnd);
            this.imageBox.PanStart += new System.EventHandler(this.imageBox_PanStart);
            this.imageBox.Paint += new System.Windows.Forms.PaintEventHandler(this.imageBox_Paint);
            this.imageBox.KeyDown += new System.Windows.Forms.KeyEventHandler(this.imageBox_KeyDown);
            this.imageBox.MouseClick += new System.Windows.Forms.MouseEventHandler(this.imageBox_MouseClick);
            // 
            // tabPageLayout
            // 
            this.tabPageLayout.Controls.Add(this.textBoxLayout);
            this.tabPageLayout.Location = new System.Drawing.Point(4, 22);
            this.tabPageLayout.Name = "tabPageLayout";
            this.tabPageLayout.Padding = new System.Windows.Forms.Padding(3);
            this.tabPageLayout.Size = new System.Drawing.Size(679, 435);
            this.tabPageLayout.TabIndex = 1;
            this.tabPageLayout.Text = "Layout";
            this.tabPageLayout.UseVisualStyleBackColor = true;
            // 
            // textBoxLayout
            // 
            this.textBoxLayout.Dock = System.Windows.Forms.DockStyle.Fill;
            this.textBoxLayout.Location = new System.Drawing.Point(3, 3);
            this.textBoxLayout.Multiline = true;
            this.textBoxLayout.Name = "textBoxLayout";
            this.textBoxLayout.ReadOnly = true;
            this.textBoxLayout.Size = new System.Drawing.Size(673, 429);
            this.textBoxLayout.TabIndex = 0;
            // 
            // treeView
            // 
            this.treeView.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.treeView.Location = new System.Drawing.Point(149, 9);
            this.treeView.Name = "treeView";
            this.treeView.Size = new System.Drawing.Size(326, 107);
            this.treeView.TabIndex = 4;
            // 
            // groupBoxCurrent
            // 
            this.groupBoxCurrent.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.groupBoxCurrent.Location = new System.Drawing.Point(481, 4);
            this.groupBoxCurrent.Name = "groupBoxCurrent";
            this.groupBoxCurrent.Size = new System.Drawing.Size(101, 112);
            this.groupBoxCurrent.TabIndex = 5;
            this.groupBoxCurrent.TabStop = false;
            this.groupBoxCurrent.Text = "Current group";
            // 
            // tabPageFcserver
            // 
            this.tabPageFcserver.Controls.Add(this.textBoxFcserver);
            this.tabPageFcserver.Location = new System.Drawing.Point(4, 22);
            this.tabPageFcserver.Name = "tabPageFcserver";
            this.tabPageFcserver.Padding = new System.Windows.Forms.Padding(3);
            this.tabPageFcserver.Size = new System.Drawing.Size(679, 435);
            this.tabPageFcserver.TabIndex = 2;
            this.tabPageFcserver.Text = "Fadecandy Config";
            this.tabPageFcserver.UseVisualStyleBackColor = true;
            // 
            // textBoxFcserver
            // 
            this.textBoxFcserver.Dock = System.Windows.Forms.DockStyle.Fill;
            this.textBoxFcserver.Location = new System.Drawing.Point(3, 3);
            this.textBoxFcserver.Multiline = true;
            this.textBoxFcserver.Name = "textBoxFcserver";
            this.textBoxFcserver.Size = new System.Drawing.Size(673, 429);
            this.textBoxFcserver.TabIndex = 0;
            // 
            // buttonSaveConfigAs
            // 
            this.buttonSaveConfigAs.Location = new System.Drawing.Point(7, 77);
            this.buttonSaveConfigAs.Name = "buttonSaveConfigAs";
            this.buttonSaveConfigAs.Size = new System.Drawing.Size(105, 23);
            this.buttonSaveConfigAs.TabIndex = 2;
            this.buttonSaveConfigAs.Text = "Save config as...";
            this.buttonSaveConfigAs.UseVisualStyleBackColor = true;
            // 
            // buttonAddFadecandy
            // 
            this.buttonAddFadecandy.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Left)));
            this.buttonAddFadecandy.Location = new System.Drawing.Point(149, 122);
            this.buttonAddFadecandy.Name = "buttonAddFadecandy";
            this.buttonAddFadecandy.Size = new System.Drawing.Size(97, 23);
            this.buttonAddFadecandy.TabIndex = 6;
            this.buttonAddFadecandy.Text = "Add Fadecandy";
            this.buttonAddFadecandy.UseVisualStyleBackColor = true;
            // 
            // buttonRemoveFadecandy
            // 
            this.buttonRemoveFadecandy.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Bottom | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonRemoveFadecandy.Location = new System.Drawing.Point(400, 122);
            this.buttonRemoveFadecandy.Name = "buttonRemoveFadecandy";
            this.buttonRemoveFadecandy.Size = new System.Drawing.Size(75, 23);
            this.buttonRemoveFadecandy.TabIndex = 7;
            this.buttonRemoveFadecandy.Text = "Remove";
            this.buttonRemoveFadecandy.UseVisualStyleBackColor = true;
            // 
            // Form
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(687, 586);
            this.Controls.Add(this.tabControl);
            this.Controls.Add(this.panel1);
            this.MinimumSize = new System.Drawing.Size(300, 300);
            this.Name = "Form";
            this.Text = "Form";
            this.panel1.ResumeLayout(false);
            this.groupBoxPoints.ResumeLayout(false);
            this.groupBoxFile.ResumeLayout(false);
            this.tabControl.ResumeLayout(false);
            this.tabPageImage.ResumeLayout(false);
            this.tabPageLayout.ResumeLayout(false);
            this.tabPageLayout.PerformLayout();
            this.tabPageFcserver.ResumeLayout(false);
            this.tabPageFcserver.PerformLayout();
            this.ResumeLayout(false);

        }

        #endregion
        private System.Windows.Forms.Panel panel1;
        private System.Windows.Forms.Button buttonOpen;
        private System.Windows.Forms.TabControl tabControl;
        private System.Windows.Forms.TabPage tabPageImage;
        private System.Windows.Forms.TabPage tabPageLayout;
        private System.Windows.Forms.TextBox textBoxLayout;
        private System.Windows.Forms.Button buttonSaveLayoutAs;
        private Cyotek.Windows.Forms.ImageBox imageBox;
        private System.Windows.Forms.GroupBox groupBoxFile;
        private System.Windows.Forms.GroupBox groupBoxPoints;
        private System.Windows.Forms.Button buttonClear;
        private System.Windows.Forms.Button buttonRemove;
        private System.Windows.Forms.TreeView treeView;
        private System.Windows.Forms.GroupBox groupBoxCurrent;
        private System.Windows.Forms.TabPage tabPageFcserver;
        private System.Windows.Forms.TextBox textBoxFcserver;
        private System.Windows.Forms.Button buttonSaveConfigAs;
        private System.Windows.Forms.Button buttonRemoveFadecandy;
        private System.Windows.Forms.Button buttonAddFadecandy;
    }
}

