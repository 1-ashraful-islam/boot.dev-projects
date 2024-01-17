from tkinter import Tk, BOTH, Canvas

class Point():
  def __init__(self, x, y):
    self.x = x
    self.y = y

class Line():
  def __init__(self, start: Point, end: Point):
    self.start = start
    self.end = end
  
  def draw(self, canvas: Canvas, fill_color):
    canvas.create_line(self.start.x, self.start.y, self.end.x, self.end.y, fill=fill_color, width = 2)
    canvas.pack()



class Window():
  def __init__(self, width, height):
    self.width = width
    self.height = height
    self.__root = Tk()
    self.__root.title("Maze Solver")
    self.__canvas = Canvas(self.__root, width=self.width, height=self.height)
    self.__canvas.pack()
    self.__is_running = False
    self.__root.protocol("WM_DELETE_WINDOW", self.close)

  def redraw(self):
    self.__root.update_idletasks()
    self.__root.update()

  def wait_for_close(self):
    self.__is_running = True
    while self.__is_running:
      self.redraw()
  
  def close(self):
    self.__is_running = False

  def draw_line(self, line: Line, fill_color = "black"):
    line.draw(self.__canvas, fill_color)

class Cell():
  def __init__(self, top_left: Point, bottom_right: Point, window: Window):
    self.__x1 = top_left.x
    self.__y1 = top_left.y
    self.__x2 = bottom_right.x
    self.__y2 = bottom_right.y
    self.has_left_wall = True
    self.has_right_wall = True
    self.has_top_wall = True
    self.has_bottom_wall = True
    self.__win = window
  
  def draw(self):
    self.__win.draw_line
    if self.has_left_wall:
      self.__win.draw_line(Line(Point(self.__x1, self.__y1), Point(self.__x1, self.__y2)))
      
    if self.has_right_wall:
      self.__win.draw_line(Line(Point(self.__x2, self.__y1), Point(self.__x2, self.__y2)))
      
    if self.has_top_wall:
      self.__win.draw_line(Line(Point(self.__x1, self.__y1), Point(self.__x2, self.__y1)))
      
    if self.has_bottom_wall:
      self.__win.draw_line(Line(Point(self.__x1, self.__y2), Point(self.__x2, self.__y2)))
  
  def draw_move(self, to_cell, undo = False):
    Point1 = Point((self.__x1 + self.__x2) / 2, (self.__y1 + self.__y2) / 2)
    Point2 = Point((to_cell.__x1 + to_cell.__x2) / 2, (to_cell.__y1 + to_cell.__y2) / 2)
    if undo:
      self.__win.draw_line(Line(Point1, Point2), "gray")
    else:
      self.__win.draw_line(Line(Point1, Point2), "red")

def main():
  window = Window(800, 600)
  cell = Cell(Point(10, 10), Point(100, 100), window)
  cell.draw()
  cell2 = Cell(Point(100, 10), Point(200, 100), window)
  cell2.draw()
  cell.draw_move(cell2)
  window.wait_for_close()

if __name__ == "__main__":
  main()
