class Window {
  constructor(width, height, background = "#ccc") {
    this.width = width;
    this.height = height;
    this.background = background;
    this.canvas = document.getElementById('myCanvas');
    this.canvas.style.background = this.background;
    this.isRunning = false;
    this.createButtons();
    this.autoSolve = true; //true by default

    // Parameters for maze generation
    
    

    this.numCols = 12;
    this.numRows = 10;


    this.padding = 50;
    this.seed = null;
    this.cellWidth = (width - 2 * this.padding) / this.numCols;
    this.cellHeight = (height - 2 * this.padding) / this.numRows;

    this.mazeProcessing = false;
    // Initial Maze creation
    this.createMaze();
  }

  createButtons() {
    const resetButton = document.getElementById('resetButton');
    resetButton.addEventListener('click', () => this.createMaze());

    const stopButton = document.getElementById('stopButton');
    stopButton.addEventListener('click', () => this.stop());

    const solveButton = document.getElementById('solveButton');
    solveButton.addEventListener('click', () => this.solveMaze(solveButton));

    const generateButton = document.getElementById('mazeGenerate');
    generateButton.addEventListener('click', () => this.generateMaze());
    
    const autoSolveCheckbox = document.getElementById('autoSolve');
    autoSolveCheckbox.addEventListener("change", () => {
      this.autoSolve = autoSolveCheckbox.checked ? true : false;
    });
  }

  generateMaze() {

    this.stop();

    const numColumnsInput = document.getElementById('numColumns');
    const numRowsInput = document.getElementById('numRows');
    this.numCols = parseInt(numColumnsInput.value) || 12; // default value of 12 if not retrieved
    this.numRows = parseInt(numRowsInput.value) || 10; // default value of 10 if not retrieved
    if (this.numCols < 2 || this.numCols > 16 || this.numRows < 2 || this.numRows > 16 || isNaN(this.numCols) || isNaN(this.numRows)) {
      alert("Invalid number of rows or columns. Use 2-16 columns and 2-16 rows.");
      return;
    }
    this.cellWidth = (this.width - 2 * this.padding) / this.numCols;
    this.cellHeight = (this.height - 2 * this.padding) / this.numRows;
    this.createMaze();
  }

  async createMaze() {
    

    if (this.mazeProcessing) {
      this.mazeProcessing = false;
      this.maze.stop();
    }
    // Clear the canvas
    this.clearCanvas();

    // disable solve button
    const solveButton = document.getElementById('solveButton');
    solveButton.disabled = true;
    solveButton.innerHTML = "Solve (wait...)";

    this.mazeProcessing = true;
    // Assuming Maze is a class you have defined in maze.js
    this.maze = new Maze(this.padding, this.padding, this.numRows, this.numCols,
                         this.cellWidth, this.cellHeight, this, this.seed);
    await this.maze.createCells();

    // enable solve button
    solveButton.disabled = false;
    solveButton.innerHTML = "Solve";
  
    // this.maze.solve();
    this.mazeProcessing = false;

    if (this.autoSolve) {
      this.solveMaze(solveButton);
    }
  }

  solveMaze(solveButton) {
    solveButton.disabled = true;
    solveButton.innerHTML = "Solving...";
    this.maze.solve();
    solveButton.disabled = false;
    solveButton.innerHTML = "Solve";
  }

  clearCanvas() {
    const ctx = this.canvas.getContext('2d');
    ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
  }

  stop() {
    if (this.mazeProcessing) {
      this.maze.stop();
    }
  }

  drawLine(line, fillColor = "black") {
    line.draw(this.canvas.getContext('2d'), fillColor);
  }

  getCanvas() {
    return this.canvas;
  }
}

// Main execution
function main() {
  const width = 800;
  const height = 600;
  const window = new Window(width, height);
  window.isRunning = true;
}

main();
