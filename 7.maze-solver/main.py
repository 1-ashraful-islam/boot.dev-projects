from tkinter import Tk, BOTH, Canvas, Button
from gprimitives import Point, Line, Cell
from maze import Maze
from time import sleep


class Window():
  def __init__(self, width, height, background="#ccc"):
    self.width = width
    self.height = height
    self.background = background
    self.__root = Tk()
    self.__root.title("Maze Solver")
    self.__canvas = Canvas(self.__root, width=self.width, height=self.height, background=self.background)
    self.__canvas.pack()
    self.__is_running = False
    self.__root.protocol("WM_DELETE_WINDOW", self.close)
    self._create_buttons()
    

    # Parameters for maze generation
    self.num_cols = 12
    self.num_rows = 10
    self.padding = 50
    self.seed = None
    self.maze_width = (width - 2 * self.padding) / self.num_cols
    self.maze_height = (height - 2 * self.padding) / self.num_rows

    self.maze_processing = False
    # Initial Maze creation
    self.create_new_maze()

  def redraw(self):
    self.__root.update_idletasks()
    self.__root.update()

  def wait_for_close(self):
    self.__is_running = True
    self.__root.mainloop()
    # while self.__is_running:
    #   print("waiting")
    #   self.redraw()
  
  def close(self):
    self.__is_running = False
    if hasattr(self, 'maze'):  # Check if maze exists before trying to stop it
        self.maze_processing = True
        if hasattr(self.maze, 'stop'):  # Check if the maze has a stop method
            self.maze.stop()
    self.__root.quit()  # Gracefully terminate the mainloop
    self.__root.destroy()  # This will close the Tkinter window immediately
    

  
  def is_running(self):
    return self.__is_running

  def draw_line(self, line: Line, fill_color = "black"):
    line.draw(self.__canvas, fill_color)

  def _create_buttons(self):
    self.__buttons = []
    self.__buttons.append(Button(self.__root, text="Exit", command=self.close))
    self.__buttons[-1].pack(side="right")
    self.__buttons.append(Button(self.__root, text="Reset", command=self.create_new_maze))
    self.__buttons[-1].pack(side="right")
    
  def create_new_maze(self):
    # Clear the canvas
    if self.maze_processing:
      self.maze_processing = False
      self.maze.stop()
    self.__canvas.delete("all")

    self.maze_processing = True
    self.maze = Maze(self.padding, self.padding, self.num_rows, self.num_cols, 
                         self.maze_width, self.maze_height, self, seed=self.seed)
    self.maze.draw()
    self.maze.solve()
    self.maze_processing = False

  def get_canvas(self):
    return self.__canvas

def main():
  width = 800
  height = 600
  window = Window(width, height)
  # num_cols = 12
  # num_rows = 10
  # padding = 50
  # seed = None
  # maze = Maze(padding, padding, num_rows, num_cols, (width - 2*padding) /num_cols, (height - 2*padding)/num_rows, window, seed= seed)
  # maze.solve()
  window.wait_for_close()

if __name__ == "__main__":
  main()
