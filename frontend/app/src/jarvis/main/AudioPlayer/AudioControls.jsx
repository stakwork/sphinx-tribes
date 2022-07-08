import React from 'react';
import { ReactComponent as Play } from '../../../utils/play-button.svg';
import { ReactComponent as Pause } from '../../../utils/pause-button.svg';
import { ReactComponent as Next } from '../../../utils/next-button.svg';
import { ReactComponent as Prev } from '../../../utils/next-button.svg';
import './audioControls.css'

const AudioControls = ({ isPlaying, onNextClick,onPrevClick, onPlayPauseClick }) => (
	<div className="audio-controls">
    <button
      type="button"
      className="prev"
      aria-label="Previous"
      onClick={onPrevClick}
      style={{transform: 'rotate(180deg)'}}
    >
      <Prev />
    </button>
    {isPlaying ? (
      <button
        type="button"
        className="pause"
        onClick={() => onPlayPauseClick(false)}
        aria-label="Pause"
      >
        <Pause />
      </button>
    ) : (
      <button
        type="button"
        className="play"
        onClick={() => onPlayPauseClick(true)}
        aria-label="Play"
      >
        <Play />
      </button>
    )}
    <button
      type="button"
      className="next"
      aria-label="Next"
      onClick={onNextClick}
    >
      <Next />
    </button>
  </div>
);

export default AudioControls;