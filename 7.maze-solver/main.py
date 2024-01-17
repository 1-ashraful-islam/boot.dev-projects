from tkinter import Tk, BOTH, Canvas

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


def main():
  window = Window(800, 600)
  window.wait_for_close()

if __name__ == "__main__":
  main()
