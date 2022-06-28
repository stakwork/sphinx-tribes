import React, {useState} from 'react'
import { EuiCard, EuiFormFieldset, EuiIcon, EuiFlexItem, EuiFlexGroup } from "@elastic/eui";
import styled from "styled-components";
import ForceGraph from './ForceGraph/ForceGraph'
//import * as d3 from "d3"

export default function BodyComponent() {

  // We should use this data for the graph
 /*const graphData = [
   {
     timestamp: "00:00:30",
     podcast: "Wat bitcoin done",
     listOfKeyWords: ["bitcoin", "taproot", "tapscript"]
   },{
     timestamp: "00:01:42",
     podcast: "yellow pill podcast",
     listOfKeyWords: ["bitcoin", "musig2"]
   },{
     timestamp: "00:00:11",
     podcast: "yeast radio podcast",
     listOfKeyWords: ["barack obama", "democrat"]
   },{
     timestamp: "00:00:11",
     podcast: "meat radio",
     listOfKeyWords: ["beef", "bitcoin", "lightning"]
   },{
     timestamp: "00:00:11",
     podcast: "meat radio",
     listOfKeyWords: ["beef", "bitcoin", "lightning"]
   },{
     timestamp: "00:00:11",
     podcast: "meat radio",
     listOfKeyWords: ["beef", "bitcoin", "lightning"]
   }
 ]
  */
  const [textBoxText, setTextBoxText] = useState("");
  const [graphData, setGraphData] = useState([])
  
  
 

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
  function callApi(word){
    fetch(`https://ardent-pastry-basement.wayscript.cloud/prediction/${word}`).then(response => response.json()).then(data => { if(data.length)
      setGraphData(data)})
    
  }
  const nodes = [
    { id: 0, r: 15, name: 'Podcast1', type: 'podcast'},
    { id: 1, r: 15, name: 'Podcast2', type: 'podcast'},
    { id: 2, r: 15, name: 'musig2', type: 'topic'},
    { id: 3, r: 15, name: 'Podcast4', type: 'podcast'},
    { id: 4, r: 15, name: 'Podcast5', type: 'podcast'},
    { id: 5, r: 15, name: 'Podcast6', type: 'podcast'},
    { id: 6, r: 15, name: 'Podcast7', type: 'podcast'},
    { id: 7, r: 15, name: 'taproot', type: 'topic'}
  ]
  const links = [
    { source: 0, target: 2},
    { source: 0, target: 7},
    { source: 3, target: 7},
    { source: 6, target: 7},
  ]
  return(
    <Body>
      <Column className="main-wrap" style={bodyStyle}>
          <form>
            <input
              style={{borderRadius: '100px', paddingLeft: '10px', marginBottom: '10px'}}
              type="text" 
              value={textBoxText} 
              placeholder="Search"
              onChange={e => {
                setTextBoxText(e.target.value)
                callApi(e.target.value)             
                }}
            />
          </form>
          {listMapping()}
          <ForceGraph linksData={links} nodesData={nodes}/>
      </Column>
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