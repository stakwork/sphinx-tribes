import React, { useState, useEffect, useRef } from 'react';
import AudioControls from './AudioControls';
import './audioPlayer.css'

const AudioPlayer = ({ tracks }) => {
  
	// State
  const [trackIndex, setTrackIndex] = useState(0);
  // Destructure for conciseness
	const { title, artist, color, image, audioSrc, timestamp } = tracks[trackIndex];
  const [trackProgress, setTrackProgress] = useState(getStartAndEndInSecond(timestamp)[0] || 0);
  const [isPlaying, setIsPlaying] = useState(true);

  

	// Refs
  const audioRef = useRef(new Audio(audioSrc));
  
  const intervalRef = useRef();
  const isReady = useRef(false);

  // Destructure for conciseness
	const { duration } = audioRef.current;
  console.log(duration)
  
  const toPrevTrack = () => {
  if (trackIndex - 1 < 0) {
    setTrackIndex(tracks.length - 1);
  } else {
    setTrackIndex(trackIndex - 1);
  }
}

const toNextTrack = () => {
  if (trackIndex < tracks.length - 1) {
    setTrackIndex(trackIndex + 1);
  } else {
    setTrackIndex(0);
  }
}

  
  useEffect(() => {
    
  if (isPlaying) {
    audioRef.current.play();
    startTimer();
  } else {
    audioRef.current.pause();
  }
}, [isPlaying]);

  useEffect(() => {
  // Pause and clean up on unmount
  return () => {
    audioRef.current.pause();
    clearInterval(intervalRef.current);
  }
}, []);

  // Handle setup when changing tracks
useEffect(() => {
  audioRef.current.pause();

  audioRef.current = new Audio(audioSrc);
  let startAndEndTime = []
  if(timestamp) 
    startAndEndTime = getStartAndEndInSecond(timestamp)
  audioRef.current.currentTime = startAndEndTime[0] || 0
	//setTrackProgress(audioRef.current.currentTime);

  if (isReady.current) {
    audioRef.current.play();
    setIsPlaying(true);
    startTimer();
  } else {
    // Set the isReady ref as true for the next pass
    isReady.current = true;
  }
}, [trackIndex]);

  const startTimer = () => {
	  // Clear any timers already running
	  clearInterval(intervalRef.current);

	  intervalRef.current = setInterval(() => {
	    if (audioRef.current.ended) {
	      toNextTrack();
	    } else {
	      setTrackProgress(audioRef.current.currentTime);
	    }
	  }, [1000]);
	}

  const onScrub = (value) => {
	// Clear any timers already running
  clearInterval(intervalRef.current);
  audioRef.current.currentTime = value;
  setTrackProgress(audioRef.current.currentTime);
}

const onScrubEnd = () => {
  // If not already playing, start
  if (!isPlaying) {
    setIsPlaying(true);
  }
  startTimer();
}

  const currentPercentage = duration ? `${(trackProgress / duration) * 100}%` : '0%';
const trackStyling = `
  -webkit-gradient(linear, 0% 0%, 100% 0%, color-stop(${currentPercentage}, #fff), color-stop(${currentPercentage}, #777))
`;

  return (
		<div className="audio-player">
			<div className="track-info">
        <h2>#{trackIndex + 1}/{tracks.length}</h2>
			  <img
			    className="artwork"
			    src={image}
			    alt={`track artwork for ${title} by ${artist}`}
			  />
		    <h2 className="title">{title}</h2>
        <h3 className="artist">{artist}</h3>
        <AudioControls          
          isPlaying={isPlaying}          
          onPrevClick={toPrevTrack}          
          onNextClick={toNextTrack}          
          onPlayPauseClick={setIsPlaying}      
          />
        <p2>{Math.trunc(trackProgress)}/{Math.trunc(duration)} seconds</p2>
        <input
        type="range"
        value={trackProgress}
        step="1"
        min="0"
        max={duration ? duration : `${duration}`}
        className="progress"
        onChange={(e) => onScrub(e.target.value)}
        onMouseUp={onScrubEnd}
        onKeyUp={onScrubEnd}
        style={{ background: trackStyling, width: '100%' }}
      />
			</div>
		</div>
	);
}

function getStartAndEndInSecond(timestamp){
  let timestampArray = timestamp.split("-")
  let timestampOneArray = timestampArray[0].split(":")
  let timestampTwoArray = timestampArray[1].split(":")

  let startTime = parseInt(timestampOneArray[0]) * 3600 + parseInt(timestampOneArray[1]) * 60 + parseInt(timestampOneArray[2])
  let endTime = parseInt(timestampTwoArray[0]) * 3600 + parseInt(timestampTwoArray[1]) * 60 + parseInt(timestampTwoArray[2])
  return [startTime, endTime]
}

export default AudioPlayer;