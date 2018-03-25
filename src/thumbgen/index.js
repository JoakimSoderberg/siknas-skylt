const fs = require('fs');
const D3Node = require('d3-node')
const d3 = require('d3');
var jsdom = require('jsdom');
const { JSDOM } = jsdom;


const markup = '<div id="container"><h3>U.S. States</h3><div id="statesgrid"></div></div>';


// Drawing to a canvas will allow us to export as PNG.
const canvasModule = require('canvas');
const d3n = new D3Node({ canvasModule: canvasModule });
//const d3nSvg = new D3Node({ canvasModule: canvasModule });
const canvas = d3n.createCanvas(960, 500);
const context = canvas.getContext('2d');

//const dom = new JSDOM(``, { runScripts: "outside-only" });

var svg = d3n.createSVG(960, 500);

var jsonCircles = [
    { "cx": 30, "cy": 30, "r": 20, "color" : "green" },
    { "cx": 70, "cy": 70, "r": 20, "color" : "purple"},
    { "cx": 110, "cy": 100, "r": 20, "color" : "red"}];

var circles = svg.selectAll("circle")
    .data(jsonCircles)
    .enter()
    .append("circle")
    .attr('cx', function(d, i) { return d['cx'] })
    .attr('cx', function(d, i) { return d['cy'] })
    .attr('r', function(d, i) { return d['r'] })
    .attr('color', function(d, i) { return d['color'] });

// draw on your canvas, then output canvas to png
canvas.pngStream().pipe(fs.createWriteStream('output.png'));

console.log(d3n.svgString());
fs.writeFile('output.svg', d3n.svgString());

//require('./lib/output')('bla', d3n);