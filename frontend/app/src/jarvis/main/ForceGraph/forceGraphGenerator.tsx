import * as d3 from 'd3';

const SELECTED_TOPIC_COLOR = '#2dc34e';
const UNSELECTED_TOPIC_COLOR = '#4a7bd8';
const PODCAST_COLOR = '#9D00A0';

export function runForceGraph(
  container: any,
  linksData: any,
  nodesData: any,
  onNodeClicked: any,
  currentTopic: string
) {
  const links = linksData.map((d: any) => Object.assign({}, d));
  const nodes = nodesData.map((d: any) => Object.assign({}, d));

  if (container.children.length > 0) {
    // TODO: Check why we need this
    d3.select('svg').remove()
  }
  const containerRect = container.getBoundingClientRect();
  const height = containerRect.height;
  const width = containerRect.width;

  const color = (d: any) => {
    if (d.name === currentTopic) return SELECTED_TOPIC_COLOR
    return d.type === 'podcast' ? PODCAST_COLOR : UNSELECTED_TOPIC_COLOR;
  };

  const drag = (simulation: any) => {
    const dragstarted = (event: any, d: any) => {
      if (!d3.active) simulation.alphaTarget(0.3).restart()
      d.fx = d.x;
      d.fy = d.y;
    };

    const dragged = (event: any, d: any) => {
      d.fx = event.x;
      d.fy = event.y;
      simulation.alpha(0.3).restart()
    };

    const dragended = (event: any, d: any) => {
      d.fx = null;
      d.fy = null;
    };

    return d3
      .drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended);
  };

  const simulation = d3
    .forceSimulation(nodes)
    .force('link', d3.forceLink(links).id((d: any) => d.id).strength(1))
    .force('charge', d3.forceManyBody().strength(-400))
    .force('x', d3.forceX())
    .force('y', d3.forceY());

  const svg = d3
    .select(container)
    .append('svg')
    .attr('viewBox', [-width / 2, -height / 2, width, height])

  const link = svg
    .append('g')
    .attr('stroke', '#999')
    .attr('stroke-opacity', 0.6)
    .selectAll('line')
    .data(links)
    .join('line')
    .attr('stroke-width', (d: any) => Math.sqrt(d.value));

  let defs = svg.append('svg:defs');
  nodes.forEach(node => {
    defs.append("svg:pattern")
      .attr("id", "grump_avatar")
      .attr("width", 15) 
      .attr("height", 15)
      .attr("patternUnits", "userSpaceOnUse")
      .append("svg:image")
      .attr("xlink:href", "https://bitcoin.org/img/icons/opengraph.png")
      .attr("width", 15)
      .attr("height", 15)
      .attr("x", 0)
      .attr("y", 0);
  })

  const node = svg
    .append('g')
    .attr('stroke', '#000')
    .attr('stroke-width', 2)
    .selectAll('circle')
    .data(nodes)
    .join('circle')
    .on('click', onNodeClicked)
    .attr('r', 15)
    .attr('fill', color)
    .on("mouseover", function(){
        // @ts-ignore
        d3.select(this)
            .style("fill", "url(#grump_avatar)")
    })
    .on("mouseout", function(){
        // @ts-ignore
        d3.select(this)
            .style("fill", color)
      
    })
    // @ts-ignore
    .call(drag(simulation));

  const label = svg.append('g')
    .attr('class', 'labels')
    .selectAll('text')
    .data(nodes)
    .enter()
    .append('text')
    .text((d: any) => d.type == 'podcast' ? null : d.name)
    .attr('text-anchor', 'middle')
    .attr('dominant-baseline', 'central')
    // @ts-ignore
    .call(drag(simulation));

  simulation.on('tick', () => {
    //update link positions
    link
      .attr('x1', (d: any) => d.source.x)
      .attr('y1', (d: any) => d.source.y)
      .attr('x2', (d: any) => d.target.x)
      .attr('y2', (d: any) => d.target.y);

    // update node positions
    node
      .attr('cx', (d: any) => d.x)
      .attr('cy', (d: any) => d.y);

    // update label positions
    label
      .attr('x', (d: any) => { return d.x; })
      .attr('y', (d: any) => { return d.y; })
  });

  function handleZoom(e) {
  d3.selectAll('g')
    .attr('transform', e.transform);
  }

  let zoom = d3.zoom()
    .on('zoom', handleZoom);
  
  d3.select('svg')
    .call(zoom);

  return {
    destroy: () => {
      d3.select('svg').remove()
      simulation.stop();
    }
  };
}
