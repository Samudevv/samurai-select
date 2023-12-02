samurai-select(1)

# NAME

samurai-select - A screen selection tool for wayland compositors using the layer shell

# SYNOPSIS

*smel* [OPTIONS]

# DESCRIPTION

samurai-select (or short smel) is a command-line utility used to select a portion of a screen of a wayland compositor using the layer shell. The geometry of the selected region will be output to standard output in the format "<x>,<y> <width>x<height>" by default (it can be configured using the *-f* flag).

Additional features are
- Taking a screenshot using the *grim* utility
- Freezing the screen if the compositor implements the wlroots screencopy protocol
- Altering the selection after performing an initial selection
- Selecting windows or regions on the screen. Following compositors are supported by default:
	- Hyprland
	- Sway

The selection can always be cancelled by pressing the _ESC_ key. This exits with code 1 and prints *selection cancelled* to standard error.

# OPTIONS

All color options support *CSS Colors* e.g. #FF67ABFF, rgb(100%, 25%, 45%), lime, red

*--background-color* _color_
	Set the color that fills the screen (default: #FFFFFF40)

*--selection-color* _color_
	Set the color that is used to draw the inside of the selection box (default: #00000000)

*--border-color* _color_
	Set the color that is used to draw the border around the selection box (default: #000000FF)

*--text-color* _color_
	Set the color that is used for the text (default: #000000FF)

*--grabber-color* _color_
	The fill color of the grabbers for altering the selection (default: #101010FF)

*--grabber-border-color* _color_
	The border color of the grabbers for altering the selection (default: #000000FF)

*--border-width* _width_
	The width of the border around the selection box in pixels (default: 2.0)

*-t*|*--text*
	Display the selection position and dimensions next to the selection box

*--font* _font family_
	Set the font family of the text (use *--list-fonts* to see which can be used) (default: sans-serif)

*--list-fonts*
	List installed fonts that can be used and exit

*--font-size* _font size_
	Set the font size of the text (default: 16)

*--text-padding* _padding_
	The distance between the selection box and the text (default: 10)

*-z*|*--freeze*
	Freeze the screen while performing the selection (wlr_screencopy protocol is required for it to work). It basically takes a screenshot and renders it beneath the surface on which the selection is rendered

*-s*|*--screenshot*
	Use *grim* to perform a screenshot. Before a screenshot is performed the part of the screen where selection is gets cleared. This is necessary if the layer surface fades away instead of disappearing instantaneously.

*-o*|*--output* _file path_
	Set the file path to where the screenshot is written. This is a format string which can be used to add information about the current date. See *FORMAT* for more information (default: screenshot-%y.%M.%d-%h:%m:%s.png)

*--screenshot-flags* _flags_
	These flags are passed to grim when performing the screenshot

*-c*|*--cmd* _command_
	Execute an arbitrary command before quitting the application. Before the command is executed the part of the screen where selection is gets cleared. This is necessary if the layer surface fades away instead of disappearing instantaneously. Insert %geometry% where you want to put the resulting geometry.

*-f*|*--format* _format string_
	Set the format of the output geometry. See *FORMAT* for more information (default: %x,%y %wx%h)

*-a*|*--aspect-ratio* _aspect ratio_
	Force an aspect ratio for the selection box in the format w:h

*-A*|*--alter-selection*
	This flag lets you change the selection box after releasing left click by dragging the box at the edges and corners. Press _Enter_ when you are done

*--grabber-radius* _radius_
	The radius of the grabbers for altering the selection (default: 7)

*--no-anim*
	Disable animations (grabber bounce when altering selection, moving of selection when selecting regions)

*-r*|*--regions* _region type_
	Choose from predefined regions of the screen. Different possible values are:
	- *auto*: The program detects which compositor is running and retrieves the window positions. This is the default value if none has been specified.
	- *hyprland*: Retrieve the window positions from Hyprland using hyprctl
	- *sway*: Retrieve the window positions from sway using swaymsg
	- *arg*: Retrive the region positions from the *-R* or *--regions-arg* flags
	- *none*: Don't select regions. This is the default one if *-r* is not used

*-R*|*--regions-arg* _regions_
	Declare a list of regions in the format 'X1,Y1 W1xH1 NAME1 X2,Y2 W2xH2 NAME2 ...'

*-p*|*--outputs*
	Select whole outputs (which is term for screens/monitors in wayland)

*-h*|*--help*
	Display a more concise help message

*-v*|*--version*
	Display version information and exit

# FORMAT

When using the *-f* or *--format* flag the following specifiers can be utilized:

%x	The absolute top left x coordinate of the selection box in pixels

%y	The absolute top left y coordinate of the selection box in pixels

%w	The width of the selection in pixels

%h	The height of the selection in pixels

%X	The output relative top left x coordinate of the selection box in pixels

%Y	The output relative top left y coordinate of the selection box in pixels

%W	The width of the selection cropped to the output

%H	The height of the selection cropped to the output

%r	The name of the region (window titles if region type *hyprland*, *sway* etc. is used)

%o	The name of the output (which is the term for screen/monitor in wayland)

The following specifiers can be used for the *-o* or *--output* flag when taking a screenshot:

%n	The nanoseconds portion of the current date

%s	The seconds portion of the current date

%m	The minutes portion of the current date

%h	The hours portion of the current date

%d	The hours portion of the current date

%M	The current month as a number

%o	The name of the current month in English

%y	The current year

# AUTHORS

This project is maintained by Jonas Pucher <jonaas.pucher000000@gmail.com> aka
PucklaJ on GitHub. See https://github.com/Samudevv/samurai-select on GitHub for the source code and ways to contribute.
