import React, {useState, useCallback, useEffect } from 'react'
import styled from "styled-components";
import ForceGraph from './ForceGraph/ForceGraph'
import * as sphinx from 'sphinx-bridge'
import AudioPlayer from './AudioPlayer/AudioPlayer'
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

interface NodesAndLinks{
  nodes: Node[],
  links: Link[]
}

const DEBOUNCE_LAG = 800

export default function BodyComponent() {
  const [initialBudget,setInitialBudget] = useState(0)
  const [topic, setTopic] = useState("");
  const [graphData, setGraphData] = useState<NodesAndLinks>({nodes: [], links: []})
  const [isLoading, setIsLoading] = useState(false)
  const [tokens, setTokens] = useState(0)
  const [pubkey,setPubkey] = useState('')
  const [validPubkey,setValidPubkey] = useState('')
  const [invoice, setInvoice] = useState({})
  const [tracks, setTracks] = useState([])
  
  // async function getOauthChallenge(){
  //   const client_id='1234567890'
  //   const q = `client_id=${client_id}&response_type=code&scope=all&mode=JSON`
  //   const url = 'https://auth.sphinx.chat/oauth?'+q
  //   try {
  //     const r1 = await fetch(url)
  //     const j = await r1.json()
  //     return j || {}
  //   } catch(e) {
  //     console.log(e)
  //     return {}
  //   }
  // }
  
  // useEffect(()=>{
  //   (async () => {
  //     const {challenge,id} = await getOauthChallenge()
  //     // @ts-ignore
  //     await sphinx.enable()
  //     // @ts-ignore
  //     const r = await sphinx.authorize(challenge, true)
  //     console.log("AUTHORIZE RES",JSON.stringify(r,null,2))
  //     if(r&&r.budget) {
  //       setInitialBudget(r.budget)
  //       setTokens(r.budget)
  //       setPubkey(r.pubkey)
  //     }
  //     if(r&&r.pubkey&&r.signature) {
  //       const r2 = await fetch(`/api/verify?id=${id}&sig=${r.signature}&pubkey=${r.pubkey}`)
  //       const j = await r2.json()
  //       console.log("VERIFY?",j)
  //       if(j&&j.valid) {
  //         setValidPubkey(j.pubkey)
  //       }
  //     }
  //   })()
  // },[])
  
  

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

  function findNodeByName(name: string, _nodes: Array<Node>) : Node | undefined {
    return _nodes.find(candidate => candidate.name === name)
  }
  
  function callApi(word: string) {
    setIsLoading(true)
    let index = 0
    // setInvoice({"value": "test"})
    // console.log(sphinx)
    // // @ts-ignore
    // let value = await sphinx.signMessage("test")
    // setInvoice(value)
    // console.log("Lookie here", value)
    fetch(`https://ardent-pastry-basement.wayscript.cloud/prediction/${word}`)
      .then(response => response.json())
      .then((data: Moment[]) => {
        if(data.length) {
          // setGraphData(data)
          const _nodes: Node[] = []
          const _links: Link[] = []
          const topicMap = {}
          // Populating nodes array with podcasts and constructing a topic map
          let tracks: any = []
          data.forEach(moment => {
            _nodes.push({
              id: index,
              name: moment.podcast_title + ":" + moment.episode_title + ":" + moment.timestamp,
              type: 'podcast'
            })
            tracks.push({
                  title: moment.podcast_title || "none",
                  artist: moment.episode_title || "none",
                  audioSrc: moment.link,
                  timestamp: moment.timestamp,
              		image: "https://thumbs.dreamstime.com/b/black-audio-wave-icon-logo-modern-sound-wave-illustration-black-audio-wave-icon-logo-modern-sound-wave-illustration-white-132130544.jpg",
                  color: "white",
            })
            index++
            const topics = moment.topics
            // @ts-ignore
            topics.forEach((topic: string) => topicMap[topic] = true)
          })
          // Adds topic nodes
          Object.keys(topicMap)
            .forEach(topic => {
              console.log("topic", topic)
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
          console.log(_nodes.filter(node => node.type == 'topic'))
          setTracks(tracks)
          setGraphData({nodes: _nodes, links: _links})
          // setNodes(_nodes)
          // console.log(_links)
          // setLinks(_links)
        }
      })
      .catch(console.error)
      .finally(() => {
        console.log('Running finally block')
        setIsLoading(false)
        console.log(isLoading)
      })
  }
  

  const dispatchNetwork = useCallback(_.debounce((word) => {
    callApi(word)
  }, DEBOUNCE_LAG), [isLoading])

  const onTopicChange = (topic: string) => {
    setTopic(topic)
    dispatchNetwork(topic)
  }

  const onNodeClicked = (event: PointerEvent, data: any, isLoading) => {
    console.log('onNodeClicked.data: ', data, ', isLoading: ', isLoading)
    if (data.type === 'topic') {
      if (!isLoading) {
        onTopicChange(data.name)
      }
    }
  }
 
  console.log(tracks.length)
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
              onChange={e => onTopicChange(e.target.value)}
            />
          </form>
      {tracks.length > 0 ? <AudioPlayer tracks={tracks}/> : null}
      {/*<ForceGraph
            linksData={graphData.links}
            nodesData={graphData.nodes}
            currentTopic={topic}
            onNodeClicked={(e,data) => onNodeClicked(e, data, isLoading)}
          />*/}
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