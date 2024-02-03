# Introduction

This maze solver is built using Python and Tkinter. It draws a randomized maze and then solves it based on the implemented algorithms. Try a live demo [here.](https://1-ashraful-islam.github.io/projects/maze-generator-and-solver/)
![mage-solver-image](https://i.imgur.com/RehzDga.png)

## Class Diagram

``` mermaid
  classDiagram
  class Point {
    +float x
    +float y
    +Point(float x, float y)
  }

  class Line {
    -Point start
    -Point end
    +Line(Point start, Point end)
    +draw(Canvas canvas, string fill_color)
  }

  class Cell {
    -float __x1
    -float __y1
    -float __x2
    -float __y2
    -Window __win
    +bool has_left_wall
    +bool has_right_wall
    +bool has_top_wall
    +bool has_bottom_wall
    +bool visited
    +Cell(Point top_left, Point bottom_right, Window window)
    +draw()
    +draw_move(Cell to_cell, bool undo)
  }

  class Window {
    -int width
    -int height
    -string background
    -Tk __root
    -Canvas __canvas
    -bool __is_running
    -Maze maze
    -int num_cols
    -int num_rows
    -int padding
    -int seed
    -float maze_width
    -float maze_height
    -bool maze_processing
    +Window(int width, int height, string background)
    +redraw()
    +wait_for_close()
    +close()
    +bool is_running()
    +draw_line(Line line, string fill_color)
    -_create_buttons()
    +create_new_maze()
  }

  class Maze {
    -int __num_cols
    -int __num_rows
    -Point __start
    -Point __end
    -bool __process
    -Cell[][] __cells
    -Window __win
    +Maze(int num_cols, int num_rows, Point start, Point end, Window window)
    +generate()
    +stop()
    +_animate()
    +_draw_emoji(string emoji, int i, int j)
    +_reset_cells_visited()
    +bool solve(string method)
    -bool _solve_r(int i, int j)
  }

  Line --|> Point : uses
  Cell --|> Point : uses
  Cell --|> Line : uses
  Cell --|> Window : uses
  Maze --|> Point : uses
  Maze --|> Cell : uses
  Maze --|> Window : uses
  ```
