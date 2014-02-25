vogl (video on OpenGL)
======================

**vogl** is a simple video player that uses OpenGL to display the frames or fields from a video stream using texture mapping and a fragment shader convert YCrCB to RGB. Several APIs are planned but the first will be a simple Go channel that sends MPEG-1 4:2:0 YCrCb pixels to the player. In that sense, **vogl** is more of a library than a standalone media player.

Possible enhancements:

• Add support for simple stereo audio using a library like PortAudio. I'll only do it if the audio and video can be in sync. I am kind of a nut for that. Nothing is worse the loss of A/V sync and so many people get it wrong or just give up.

• Support for other pixel formats including MPEG-2  


*However, the overall goal is to keep **vogl** very simple*

***Currently under heavy development***
