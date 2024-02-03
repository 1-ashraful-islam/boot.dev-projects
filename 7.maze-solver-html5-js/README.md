# Introduction

This maze solver is built using Python and Tkinter. It draws a randomized maze and then solves it based on the implemented algorithms. Try a live demo [here.](https://1-ashraful-islam.github.io/projects/maze-generator-and-solver/)
![mage-solver-image](https://i.imgur.com/RehzDga.png)

## Class Diagram

``` mermaid
  classDiagram
      class Maze {
        +int x1
        +int y1
        +int numRows
        +int numCols
        +int cellSizeX
        +int cellSizeY
        +Window window
        +seed
        +boolean process
        +Array cells
        +Array start
        +Array end
        +String savedCanvas
        +constructor()
        +createCells()
        +restoreCanvas()
        +drawAllCells()
        +drawCell()
        +stop()
        +animate()
        +breakEntranceAndExit()
        +updateWall()
        +drawCellAndEmoji()
        +breakWalls()
        +breakWallsRecursive()
        +resetCellsVisited()
        +solve()
        +solveRecursive()
      }

      class MazeWindow {
        +int width
        +int height
        +String background
        +Canvas canvas
        +CanvasRenderingContext2D ctx
        +boolean isRunning
        +boolean autoSolve
        +boolean mazeProcessing
        +int padding
        +seed
        +int numCols
        +int numRows
        +int cellWidth
        +int cellHeight
        +Maze maze
        +constructor()
        +calculate_maze_parameters()
        +createButtons()
        +generateMaze()
        +createMaze()
        +solveMaze()
        +clearCanvas()
        +stop()
        +drawLine()
        +getCanvas()
      }

      class Point {
        +int x
        +int y
        +constructor(x, y)
      }

      class Line {
        +Point start
        +Point end
        +constructor(start, end)
        +draw(ctx, fillColor, lineWidth)
      }

      class Cell {
        +int _x1
        +int _y1
        +int _x2
        +int _y2
        +boolean hasLeftWall
        +boolean hasRightWall
        +boolean hasTopWall
        +boolean hasBottomWall
        +Canvas canvas
        +CanvasRenderingContext2D ctx
        +boolean visited
        +constructor(topLeft, bottomRight, canvas)
        +draw()
        +drawMove(toCell, undo)
        +min(a, b)
        +drawEmoji(emoji, fontSize)
      }

      MazeWindow --> Maze: contains
      Cell --> Line: uses
      Cell --> Point: uses
  ```
