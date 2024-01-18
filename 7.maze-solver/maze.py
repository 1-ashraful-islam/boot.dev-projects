from gprimitives import Point, Cell
from time import sleep
import random

class Maze():
  def __init__(
          self,
          x1,
          y1,
          num_rows,
          num_cols,
          cell_size_x,
          cell_size_y,
          window,
          seed = None,
        ):
    self.__x1 = x1
    self.__y1 = y1
    self.__num_rows = num_rows
    self.__num_cols = num_cols
    self.__cell_size_x = cell_size_x
    self.__cell_size_y = cell_size_y
    self.__win = window
    self.__seed = seed
    self.__process = True

    if self.__seed:
      random.seed(self.__seed)
    self.__cells = []
    

  def _create_cells(self):
    for row in range(self.__num_rows):
      cells = []
      for col in range(self.__num_cols):
        if not self.__process:
          return
        top_left = Point(self.__x1 + col * self.__cell_size_x, self.__y1 + row * self.__cell_size_y)
        bottom_right = Point(top_left.x + self.__cell_size_x, top_left.y + self.__cell_size_y)
        cell = Cell(top_left, bottom_right, self.__win)
        cells.append(cell)
      self.__cells.append(cells) # top level of list is row, second level is column
    
    for row in range(self.__num_rows):
      for col in range(self.__num_cols):
        if not self.__process:
          return
        self._draw_cell(row, col)
        
  def draw(self):
    self._create_cells()
    self._break_entrance_and_exit()
    self._break_walls_r(self.__start[0],self.__start[1])
    self._reset_cells_visited()
    
    
  
  def _draw_cell(self, row, col, duration = 0.025):
    if self.__process:
      self.__cells[row][col].draw()
      self._animate(duration=duration)

  def stop(self):
    self.__process = False
  
  def _animate(self, duration = 0.05):
    if self.__process:
      self.__win.redraw()
      sleep(duration)

  def _break_entrance_and_exit(self):
    # select a random cell on any of the 4 sides to be the entrance
    i_start = random.randint(0, self.__num_rows - 1)
    if i_start == 0 or i_start == self.__num_rows - 1:
      j_start = random.randint(0, self.__num_cols - 1)
    else:
      j_start = random.choice([0, self.__num_cols - 1])

    i_end = i_start
    j_end = j_start
    while (abs(i_start - i_end) + abs(j_start - j_end) < (self.__num_rows + self.__num_cols)/2) and self.__process:
      i_end = random.randint(0, self.__num_rows - 1)
      if i_end == 0 or i_end == self.__num_rows - 1:
        j_end = random.randint(0, self.__num_cols - 1)
      else:
        j_end = random.choice([0, self.__num_cols - 1])

    self.__start = [i_start, j_start]
    self.__end = [i_end, j_end]

    if i_start == 0:
      self.__cells[i_start][j_start].has_top_wall = False
    elif i_start == self.__num_rows - 1:
      self.__cells[i_start][j_start].has_bottom_wall = False
    elif j_start == 0:
      self.__cells[i_start][j_start].has_left_wall = False
    elif j_start == self.__num_cols - 1:
      self.__cells[i_start][j_start].has_right_wall = False
    
    

    if i_end == 0:
      self.__cells[i_end][j_end].has_top_wall = False
    elif i_end == self.__num_rows - 1:
      self.__cells[i_end][j_end].has_bottom_wall = False
    elif j_end == 0:
      self.__cells[i_end][j_end].has_left_wall = False
    elif j_end == self.__num_cols - 1:
      self.__cells[i_end][j_end].has_right_wall = False

    if self.__process:
      self._draw_cell(i_start, j_start)
      self.__cells[i_start][j_start]._draw_emoji("🚀")

      self._draw_cell(i_end, j_end)
      self.__cells[i_end][j_end]._draw_emoji("🌔")

  def _break_walls(self, i, j, i_next, j_next):
    if i_next > i:
      self.__cells[i][j].has_bottom_wall = False
      self.__cells[i_next][j_next].has_top_wall = False
    elif i_next < i:
      self.__cells[i][j].has_top_wall = False
      self.__cells[i_next][j_next].has_bottom_wall = False
    elif j_next > j:
      self.__cells[i][j].has_right_wall = False
      self.__cells[i_next][j_next].has_left_wall = False
    elif j_next < j:
      self.__cells[i][j].has_left_wall = False
      self.__cells[i_next][j_next].has_right_wall = False
    
    self._draw_cell(i, j)
    self._draw_cell(i_next, j_next)
  

  def _break_walls_r(self, i, j):
    self.__cells[i][j].visited = True

    while self.__process:
      neighbours = []
      if i > 0 and not self.__cells[i-1][j].visited:
        neighbours.append((i-1, j))
      if i < self.__num_rows - 1 and not self.__cells[i+1][j].visited:
        neighbours.append((i+1, j))
      if j > 0 and not self.__cells[i][j-1].visited:
        neighbours.append((i, j-1))
      if j < self.__num_cols - 1 and not self.__cells[i][j+1].visited:
        neighbours.append((i, j+1))
      
      if len(neighbours) == 0:
        break
      
      (i_next, j_next) = random.choice(neighbours)
      self._break_walls(i, j, i_next, j_next)
      self._break_walls_r(i_next, j_next)

  def _reset_cells_visited(self):
    for row in range(self.__num_rows):
      for col in range(self.__num_cols):
        if not self.__process:
          return
        self.__cells[row][col].visited = False


  def solve(self, method = "r"):

    return self._solve_r(self.__start[0], self.__start[1])
  
  def _solve_r(self, i, j):
    
    self.__cells[i][j].visited = True
    self._animate()
    
    if (i == self.__end[0] and j == self.__end[1]) or not self.__process:
      return True
    
    dead_end = 0
    # depth first search with a random direction to choose from
    
    # down
    if i < self.__num_rows - 1 and not self.__cells[i+1][j].visited and not self.__cells[i][j].has_bottom_wall:
      self.__cells[i][j].draw_move(self.__cells[i+1][j])
      if self._solve_r(i+1, j):
        return True
      else:
        self.__cells[i][j].draw_move(self.__cells[i+1][j], undo = True)
    else:
      dead_end += 1
    

    # right
    if j < self.__num_cols - 1 and not self.__cells[i][j+1].visited and not self.__cells[i][j].has_right_wall:
      self.__cells[i][j].draw_move(self.__cells[i][j+1])
      if self._solve_r(i, j+1):
        return True
      else:
        self.__cells[i][j].draw_move(self.__cells[i][j+1], undo = True)
    else:
      dead_end += 1

    
    # left
    if j > 0 and not self.__cells[i][j-1].visited and not self.__cells[i][j-1].has_right_wall:
      self.__cells[i][j].draw_move(self.__cells[i][j-1])
      if self._solve_r(i, j-1):
        return True
      else:
        self.__cells[i][j].draw_move(self.__cells[i][j-1], undo = True)
    else:
      dead_end += 1

    # up
    if i > 0 and not self.__cells[i-1][j].visited and not self.__cells[i-1][j].has_bottom_wall:
      self.__cells[i][j].draw_move(self.__cells[i-1][j])
      if self._solve_r(i-1, j):
        return True
      else:
        self.__cells[i][j].draw_move(self.__cells[i-1][j], undo = True)
    else:
      dead_end += 1
    
    if dead_end == 4 and self.__process:
      self.__cells[i][j]._draw_emoji("💥")
     
