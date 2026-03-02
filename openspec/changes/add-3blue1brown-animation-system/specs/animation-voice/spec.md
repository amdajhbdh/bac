## ADDED Requirements

### Requirement: Voice narration is optional
Voice narration SHALL be optional - animations work without audio.

#### Scenario: Generate animation without audio
- **WHEN** include_audio is false or not specified
- **THEN** the render SHALL proceed without any TTS processing

#### Scenario: Generate animation with audio
- **WHEN** include_audio is true
- **THEN** the system SHALL generate narration audio and merge with video

### Requirement: Narration generator creates scripts
The voice system SHALL generate narration scripts from problem solutions.

#### Scenario: Generate French narration script
- **WHEN** language is "fr" and problem is provided
- **THEN** the script builder SHALL create French narration text

#### Scenario: Generate English narration script
- **WHEN** language is "en" and problem is provided
- **THEN** the script builder SHALL create English narration text

#### Scenario: Include step-by-step narration
- **WHEN** solution contains multiple steps
- **THEN** the script SHALL include narration for each step

### Requirement: TTS synthesizes speech
The voice system SHALL use Coqui TTS to synthesize speech from narration scripts.

#### Scenario: Synthesize speech
- **WHEN** narration script is provided
- **THEN** TTS SHALL produce WAV audio file

#### Scenario: TTS failure handling
- **WHEN** TTS fails to synthesize
- **THEN** the system SHALL log warning and continue without audio

### Requirement: Audio-video synchronization
The voice system SHALL synchronize audio with video timeline.

#### Scenario: Calculate audio timing
- **WHEN** narration script has multiple timing markers
- **THEN** audio segments SHALL be timed to match video progression

#### Scenario: Merge audio with video
- **WHEN** audio file exists
- **THEN** FFmpeg SHALL merge audio with video stream
