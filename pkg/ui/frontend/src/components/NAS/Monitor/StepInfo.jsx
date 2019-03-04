import React from 'react';
import { withStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';

import * as d3 from "d3";



const styles = theme => ({
})

let range = n => [...Array(n).keys()];

let rand = (min, max) => Math.random() * (max - min) + min;

Array.prototype.last = function() { return this[this.length - 1]; };

let flatten = (array) => array.reduce((flat, toFlatten) => (flat.concat(Array.isArray(toFlatten) ? flatten(toFlatten) : toFlatten)), []);


class StepInfo extends React.Component {

    componentDidMount() {
        const id = `#graph-container${this.props.step.id}`
        let svg = d3.select(id).append("svg").attr("xmlns", "http://www.w3.org/2000/svg");
        let g = svg.append("g");
        svg.style("cursor", "move");
        let color1 = '#e0e0e0';
        let color2 = '#a0a0a0';
        let borderWidth = 1.0;
        let borderColor = "black";
        let rectOpacity = 0.8;
        let betweenSquares = 2;
        let betweenLayersDefault = 4;
        // educated guess
        let w = 150;
        let h = 360;

        
        let architecture = [];
        let lenet = {};
        let layer_offsets = [];
        let largest_layer_width = 0;
        let showLabels = true;

        let textFn = (layer) => (typeof(layer) === "object" ? layer['numberOfSquares']+'@'+layer['squareHeight']+'x'+layer['squareWidth'] : "1x"+layer)

        let rect, conv, link, poly, line, text, info;

        let architectureCNN = [
            {
                "numberOfSquares": 8,
                "squareHeight": 128,
                "squareWidth": 128,
                "filterHeight": 8,
                "filterWidth": 8,
                "op": "Max-Pool",
                "layer": 0
            },
            {
                "numberOfSquares": 8,
                "squareHeight": 64,
                "squareWidth": 64,
                "filterHeight": 16,
                "filterWidth": 16,
                "op": "Convolution",
                "layer": 1
            },
            {
                "numberOfSquares": 24,
                "squareHeight": 48,
                "squareWidth": 48,
                "filterHeight": 8,
                "filterWidth": 8,
                "op": "Max-Pool",
                "layer": 2
            },
            {
                "numberOfSquares": 24,
                "squareHeight": 16,
                "squareWidth": 16,
                "filterHeight": 8,
                "filterWidth": 8,
                "op": "Dense",
                "layer": 3
            }
        ]
        let architectureFFN = [256, 128]
        let betweenLayers = [40, 10, -20, -20]

        lenet.rects = architectureCNN.map((layer, layer_index) => range(layer['numberOfSquares']).map(rect_index => {return {'id':layer_index+'_'+rect_index,'layer':layer_index,'rect_index':rect_index,'width':layer['squareWidth'],'height':layer['squareHeight']}}));
        lenet.rects = flatten(lenet.rects);

        lenet.convs = architectureCNN.map((layer, layer_index) => Object.assign({'id':'conv_'+layer_index,'layer':layer_index}, layer)); lenet.convs.pop();
        lenet.convs = lenet.convs.map(conv => Object.assign({'x_rel':rand(0.1, 0.9),'y_rel':rand(0.1, 0.9)}, conv))

        lenet.conv_links = lenet.convs.map(conv => {return [Object.assign({'id':'link_'+conv['layer']+'_0','i':0},conv), Object.assign({'id':'link_'+conv['layer']+'_1','i':1},conv)]});
        lenet.conv_links = flatten(lenet.conv_links);

        lenet.fc_layers = architectureFFN.map((size, fc_layer_index) => {return {'id': 'fc_'+fc_layer_index, 'layer':fc_layer_index+architectureCNN.length, 'size':size/Math.sqrt(2)}});
        lenet.fc_links = lenet.fc_layers.map(fc => { return [Object.assign({'id':'link_'+fc['layer']+'_0','i':0,'prevSize':10},fc), Object.assign({'id':'link_'+fc['layer']+'_1','i':1,'prevSize':10},fc)]});
        lenet.fc_links = flatten(lenet.fc_links);
        lenet.fc_links[0]['prevSize'] = 0;                            // hacks
        lenet.fc_links[1]['prevSize'] = lenet.rects.last()['width'];  // hacks

        let label = architectureCNN.map((layer, layer_index) => { return {'id':'data_'+layer_index+'_label','layer':layer_index,'text':textFn(layer)}})
                             .concat(architectureFFN.map((layer, layer_index) => { return {'id':'data_'+layer_index+architectureCNN.length+'_label','layer':layer_index+architectureCNN.length,'text':textFn(layer)}}) );


        g.selectAll('*').remove();

        rect = g.selectAll(".rect")
                .data(lenet.rects)
                .enter()
                .append("rect")
                .attr("class", "rect")
                .attr("id", d => d.id)
                .attr("width", d => d.width)
                .attr("height", d => d.height);

        conv = g.selectAll(".conv")
                .data(lenet.convs)
                .enter()
                .append("rect")
                .attr("class", "conv")
                .attr("id", d => d.id)
                .attr("width", d => d.filterWidth)
                .attr("height", d => d.filterHeight)
                .style("fill-opacity", 0);

        link = g.selectAll(".link")
                .data(lenet.conv_links)
                .enter()
                .append("line")
                .attr("class", "link")
                .attr("id", d => d.id);

        poly = g.selectAll(".poly")
                .data(lenet.fc_layers)
                .enter()
                .append("polygon")
                .attr("class", "poly")
                .attr("id", d => d.id);

        line = g.selectAll(".line")
                .data(lenet.fc_links)
                .enter()
                .append("line")
                .attr("class", "line")
                .attr("id", d => d.id);

        text = g.selectAll(".text")
                .data(architecture)
                .enter()
                .append("text")
                .text(d => (showLabels ? d.op : ""))
                .attr("class", "text")
                .attr("dy", ".35em")
                .style("font-size", "16px")
                .attr("font-family", "sans-serif");

        info = g.selectAll(".info")
                .data(label)
                .enter()
                .append("text")
                .text(d => (showLabels ? d.text : ""))
                .attr("class", "info")
                .attr("dy", "-0.3em")
                .style("font-size", "16px")
                .attr("font-family", "sans-serif");

                rect.style("fill", d => d.rect_index % 2 ? color1 : color2);
        poly.style("fill", color1);

        rect.style("stroke", borderColor);
        conv.style("stroke", borderColor);
        link.style("stroke", borderColor);
        poly.style("stroke", borderColor);
        line.style("stroke", borderColor);

        rect.style("stroke-width", borderWidth);
        conv.style("stroke-width", borderWidth);
        link.style("stroke-width", borderWidth / 2);
        poly.style("stroke-width", borderWidth);
        line.style("stroke-width", borderWidth / 2);

        rect.style("opacity", rectOpacity);
        conv.style("stroke-opacity", rectOpacity);
        link.style("stroke-opacity", rectOpacity);
        poly.style("opacity", rectOpacity);
        line.style("stroke-opacity", rectOpacity);

        text.text(d => (showLabels ? d.op : ""));
        info.text(d => (showLabels ? d.text : ""));

        let layer_widths = architectureCNN.map((layer, i) => (layer['numberOfSquares']-1) * betweenSquares + layer['squareWidth']);
        layer_widths = layer_widths.concat(lenet.fc_layers.map((layer, i) => layer['size']));

        largest_layer_width = Math.max(...layer_widths);

        let layer_x_offsets = layer_widths.reduce((offsets, layer_width, i) => offsets.concat([offsets.last() + layer_width + (betweenLayers[i] || betweenLayersDefault) ]), [0]);
        let layer_y_offsets = layer_widths.map(layer_width => (largest_layer_width - layer_width) / 2);

        let screen_center_x = w/2 - architecture.length * largest_layer_width/2;
        let screen_center_y = h/2 - largest_layer_width/2;

        let x = (layer, node_index) => layer_x_offsets[layer] + (node_index * betweenSquares) + screen_center_x;
        let y = (layer, node_index) => layer_y_offsets[layer] + (node_index * betweenSquares) + screen_center_y;

        rect.attr('x', d => x(d.layer, d.rect_index))
            .attr('y', d => y(d.layer, d.rect_index));

        let xc = (d) => (layer_x_offsets[d.layer]) + ((d['numberOfSquares']-1) * betweenSquares) + (d['x_rel'] * (d['squareWidth'] - d['filterWidth'])) + screen_center_x;
        let yc = (d) => (layer_y_offsets[d.layer]) + ((d['numberOfSquares']-1) * betweenSquares) + (d['y_rel'] * (d['squareHeight'] - d['filterHeight'])) + screen_center_y;

        conv.attr('x', d => xc(d))
            .attr('y', d => yc(d));

        link.attr("x1", d => xc(d) + d['filterWidth'])
            .attr("y1", d => yc(d) + (d.i ? 0 : d['filterHeight']))
            .attr("x2", d => (layer_x_offsets[d.layer+1]) + ((architectureCNN[d.layer+1]['numberOfSquares']-1) * betweenSquares) + architectureCNN[d.layer+1]['squareWidth'] * d.x_rel + screen_center_x)
            .attr("y2", d => (layer_y_offsets[d.layer+1]) + ((architectureCNN[d.layer+1]['numberOfSquares']-1) * betweenSquares) + architectureCNN[d.layer+1]['squareHeight'] * d.y_rel + screen_center_y);


        poly.attr("points", function(d) {
            return ((layer_x_offsets[d.layer]+screen_center_x)           +','+(layer_y_offsets[d.layer]+screen_center_y)+
                ' '+(layer_x_offsets[d.layer]+screen_center_x+10)        +','+(layer_y_offsets[d.layer]+screen_center_y)+
                ' '+(layer_x_offsets[d.layer]+screen_center_x+d.size+10) +','+(layer_y_offsets[d.layer]+screen_center_y+d.size)+
                ' '+(layer_x_offsets[d.layer]+screen_center_x+d.size)    +','+(layer_y_offsets[d.layer]+screen_center_y+d.size));
        });

        line.attr("x1", d => layer_x_offsets[d.layer-1] + (d.i ? 0 : layer_widths[d.layer-1]) + d.prevSize + screen_center_x)
            .attr("y1", d => layer_y_offsets[d.layer-1] + (d.i ? 0 : layer_widths[d.layer-1]) + screen_center_y)
            .attr("x2", d => layer_x_offsets[d.layer] + (d.i ? 0 : d.size) + screen_center_x)
            .attr("y2", d => layer_y_offsets[d.layer] + (d.i ? 0 : d.size) + screen_center_y);

        text.attr('x', d => (layer_x_offsets[d.layer] + layer_widths[d.layer] + layer_x_offsets[d.layer+1] + layer_widths[d.layer+1]/2)/2 + screen_center_x -15)
            .attr('y', d => layer_y_offsets[0] + screen_center_y + largest_layer_width);

        info.attr('x', d => layer_x_offsets[d.layer] + screen_center_x)
            .attr('y', d => layer_y_offsets[d.layer] + screen_center_y - 15);
            svg.attr("width", 1024).attr("height", 320);
    }

    render () {
        const { step, classes } = this.props;
        const id = `graph-container${this.props.step.id}`
        return (
            <div>
                <Typography variant={"h6"}>
                    Architecture
                </Typography> 
                <div id={id}/>
                <br />
                <Typography variant={"h6"}>
                    {step.metricsName}: {step.metricsValue}.
                </Typography>
                <br />
                <a href={`${step.link}`}>
                    <Button variant={"contained"} color={"primary"}>
                        Download
                    </Button>
                </a>
            </div>
        )
    }
}

export default withStyles(styles)(StepInfo);