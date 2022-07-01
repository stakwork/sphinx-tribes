import React, {useState, useCallback } from 'react'
import { EuiCard, EuiFormFieldset, EuiIcon, EuiFlexItem, EuiFlexGroup } from "@elastic/eui";
import styled from "styled-components";
import ForceGraph from './ForceGraph/ForceGraph'
import _ from 'lodash'
import './body.css'

interface Node {
  id: number,
  name: string,
  type: string
}

interface Link {
  source: number,
  target: number
}

interface Moment {
  episode_title: string,
  link: string,
  podcast_title: string,
  timestamp: string,
  topics: string[]
}

const DEBOUNCE_LAG = 800

export default function BodyComponent() {

  const [topic, setTopic] = useState("");
  const [graphData, setGraphData] = useState([])
  const [nodes, setNodes] = useState<Node[]>([])
  const [links, setLinks] = useState<Link[]>([])
  const [isLoading, setIsLoading] = useState(false)

 const hStyle = {
    color: 'white',
   
    fontFamily: 'Roboto',
    fontWeight: 'lighter',
    fontSize: 30,
  };

  const bodyStyle = {
    height: 1000,
    background: '#212529',
  };

  const divStyle = {
    height: 10
  };

  function LineItemComponent(props){
    
    return(
      <EuiFlexItem grow={false} key={props.index}>
        <EuiCard
          icon={<EuiIcon size="xxl" type={`https://upload.wikimedia.org/wikipedia/commons/0/02/SVG_logo.svg`} />}
          title={`test`}
          isDisabled={true}
          description={props.lineItemData['0']}
          onClick={() => {}}
        />
      </EuiFlexItem >
    )
  }

  function listMapping(){
    return (
        <div style={{display:'flex', gap: '10px', flexWrap: 'wrap'}}>
          {graphData.length && graphData.map((item,index) => {

               return <LineItemComponent key={index} lineItemData={item} index/>
                 
                })}
        </div >
    )
  }

  function findNodeByName(name: string, _nodes: Array<Node>) : Node | undefined {
    return _nodes.find(candidate => candidate.name === name)
  }
  
  function callApi(word: string) {
    setIsLoading(true)
    let index = 0
    fetch(`https://ardent-pastry-basement.wayscript.cloud/prediction/${word}`)
      .then(response => response.json())
      .then((data: Moment[]) => {
        if(data.length) {
          // setGraphData(data)
          const _nodes: Node[] = []
          const _links: Link[] = []
          const topicMap = {}
          // Populating nodes array with podcasts and constructing a topic map
          data.forEach(moment => {
            _nodes.push({
              id: index,
              name: moment.podcast_title + ":" + moment.episode_title + ":" + moment.timestamp,
              type: 'podcast'
            })
            index++
            const topics = moment.topics
            // @ts-ignore
            topics.forEach((topic: string) => topicMap[topic] = true)
          })
          // Adds topic nodes
          Object.keys(topicMap)
            .forEach(topic => {
              const topicNode: Node = {
                id: index,
                name: topic,
                type: 'topic'
              }
              _nodes.push(topicNode)
              index++
            })
          // Populating the links array next
          data.forEach(moment => {
            const { topics } = moment
            topics.forEach(topic => {
              const podcastNode = findNodeByName(moment.podcast_title + ":" + moment.episode_title + ":" + moment.timestamp, _nodes)
              const topicNode = findNodeByName(topic, _nodes)
              if (podcastNode && topicNode) {
                const link: Link = {
                  source: podcastNode.id,
                  target: topicNode.id
                }
                _links.push(link)
              }
            })
          })
          setNodes(_nodes)
          setLinks(_links)
        }
      })
      .catch(console.error)
      .finally(() => setIsLoading(false))
  }

  const dispatchNetwork = useCallback(_.debounce((word) => {
    callApi(word)
  }, DEBOUNCE_LAG), [])

  const onChange = (event: any) => {
    setTopic(event.target.value)
    dispatchNetwork(event.target.value)
  }

  const onNodeClicked = (event: PointerEvent, data: any) => {
    if (data.type === 'topic') {
      if (!isLoading) {
        setTopic(data.name)
        callApi(data.name)
      }
    }
  }
  
  return(
    <Body>
          <form>
            <input
              className={isLoading ? 'loading' : ''}
              disabled={isLoading}
              style={{borderRadius: '100px', paddingLeft: '10px', marginBottom: '10px'}}
              type="text" 
              value={topic}
              placeholder="Search"
              onChange={onChange}
            />
          </form>
          {/*listMapping()*/}
          <ForceGraph
            linksData={links}
            nodesData={nodes}
            currentTopic={topic}
            onNodeClicked={onNodeClicked}
          />
    </Body>
  )
}

const Body = styled.div`
  flex:1;
  height:calc(100vh - 60px);
  // padding-bottom:80px;
  width:100%;
  overflow:auto;
  background:#272c4b;
  display:flex;
  flex-direction:column;
  align-items:center;
`

const Column = styled.div`
  display:flex;
  flex-direction:column;
  align-items:center;
  margin-top:10px;
  // max-width:900px;
  width:100%;
`