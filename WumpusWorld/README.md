### Wumpus World

The neighborhood of a node consists of the four squares north, south, east, west of the given square.
In a square the agent gets a vector of percepts, with components

Stench,Breeze,Glitter,Bump,Scream
For example [Stench,None,Glitter,None,None]
Stench is perceived at a square iff the wumpus is at this square or in its neighborhood.
Breeze is perceived at a square iff a pit is in the neighborhood of this square.
Glitter is perceived at a square iff gold is in this square
Bump is perceived at a square iff the agent goes Forward into a wall
Scream is perceived at a square iff the wumpus is killed anywhere in the cave
An agent can do the following actions (one at a time):

Turn(Right), Turn(Left), Forward, Shoot, Grab, Release, Climb
The agent can go Forward in the direction it is currently facing, or Turn Right, or Turn Left. Going Forward into a wall will generate a Bump percept.
The agent has a single arrow that it can Shoot. It will go straight in the direction faced by the agent until it hits (and kills) the wumpus, or hits (and is absorbed by) a wall.
The agent can Grab a portable object at the current square or it can Release an object that it is holding.
The agent can Climb out of the cave if at the Start square.
The Start square is (1,1) and initially the agent is facing east. The agent dies if it is in the same square as the wumpus.
The objective of the game is to kill the wumpus, to pick up the gold, and to climb out with it.
