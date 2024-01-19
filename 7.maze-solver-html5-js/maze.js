class Maze {
  constructor(x1, y1, numRows, numCols, cellSizeX, cellSizeY, window, seed = null) {
    this.x1 = x1;
    this.y1 = y1;
    this.numRows = numRows;
    this.numCols = numCols;
    this.cellSizeX = cellSizeX;
    this.cellSizeY = cellSizeY;
    this.window = window;
    this.seed = seed;
    this.process = true;

    if (this.seed) {
      Math.seedrandom(this.seed); // Requires seedrandom library
    }
    this.cells = [];
    this.start = null;
    this.end = null;

  }

  async createCells() {
    for (let row = 0; row < this.numRows; row++) {
      let cellsRow = [];
      for (let col = 0; col < this.numCols; col++) {
        let topLeft = new Point(this.x1 + col * this.cellSizeX, this.y1 + row * this.cellSizeY);
        let bottomRight = new Point(topLeft.x + this.cellSizeX, topLeft.y + this.cellSizeY);
        let cell = new Cell(topLeft, bottomRight, this.window.getCanvas());
        cellsRow.push(cell);
      }
      this.cells.push(cellsRow);
    }

    await this.drawAllCells();
    await this.breakEntranceAndExit();
    await this.breakWallsRecursive(this.start[0], this.start[1]);
    this.resetCellsVisited();

    //save the canvas to restore later before solving
    this.savedCanvas = this.window.getCanvas().toDataURL();
    
  }

  restoreCanvas() {
    const img = new Image();
    img.onload = () => {
        const ctx = this.window.getCanvas().getContext('2d');
        ctx.clearRect(0, 0, this.window.getCanvas().width, this.window.getCanvas().height);
        ctx.drawImage(img, 0, 0);
    };
    img.src = this.savedCanvas;
  }

  async drawAllCells() {
    for (let row = 0; row < this.numRows; row++) {
      for (let col = 0; col < this.numCols; col++) {
        if (this.process) {
          await this.drawCell(row, col, 30);
        }
         
      }
    }
  }

  async drawCell(row, col, duration = 0) {
    if (this.process) {
      this.cells[row][col].draw();
      await this.animate(duration);
    }
  }

  stop() {
    this.process = false;
  }

  async animate(duration = 500) {
    if (this.process) {
      await new Promise(resolve => setTimeout(resolve, duration));
    }
  }

  async breakEntranceAndExit() {
    // Select a random cell on any of the 4 sides to be the entrance
    let iStart = Math.floor(Math.random() * this.numRows);
    let jStart;
    if (iStart === 0 || iStart === this.numRows - 1) {
      jStart = Math.floor(Math.random() * this.numCols);
    } else {
      jStart = Math.random() < 0.5 ? 0 : this.numCols - 1;
    }

    let iEnd = iStart;
    let jEnd = jStart;
    let tries = 0;
    while ((Math.abs(iStart - iEnd) + Math.abs(jStart - jEnd) < (this.numRows + this.numCols) / 2) && this.process && tries++ < 100) {
      iEnd = Math.floor(Math.random() * this.numRows);
      if (iEnd === 0 || iEnd === this.numRows - 1) {
        jEnd = Math.floor(Math.random() * this.numCols);
      } else {
        jEnd = Math.random() < 0.5 ? 0 : this.numCols - 1;
      }

    }

    this.start = [iStart, jStart];
    this.end = [iEnd, jEnd];

    this.updateWall(iStart, jStart, iEnd, jEnd);
    if (this.process) {
      this.drawCellAndEmoji(iStart, jStart, "🚀");
      await this.animate(100);
      this.drawCellAndEmoji(iEnd, jEnd, "🌔");
    }
  }

  updateWall(iStart, jStart, iEnd, jEnd) {
    // Update the walls for the entrance
    if (iStart === 0) {
      this.cells[iStart][jStart].hasTopWall = false;
    } else if (iStart === this.numRows - 1) {
      this.cells[iStart][jStart].hasBottomWall = false;
    } else if (jStart === 0) {
      this.cells[iStart][jStart].hasLeftWall = false;
    } else if (jStart === this.numCols - 1) {
      this.cells[iStart][jStart].hasRightWall = false;
    }

    // Update the walls for the exit
    if (iEnd === 0) {
      this.cells[iEnd][jEnd].hasTopWall = false;
    } else if (iEnd === this.numRows - 1) {
      this.cells[iEnd][jEnd].hasBottomWall = false;
    } else if (jEnd === 0) {
      this.cells[iEnd][jEnd].hasLeftWall = false;
    } else if (jEnd === this.numCols - 1) {
      this.cells[iEnd][jEnd].hasRightWall = false;
    }
  }

  drawCellAndEmoji(i, j, emoji) {
    if (this.process) {
      this.cells[i][j].draw(); // Assuming Cell class has a draw method
      this.cells[i][j].drawEmoji(emoji); // Assuming Cell class has a drawEmoji method
    }
  }

  async breakWalls(i, j, iNext, jNext) {
    if (iNext > i) {
      this.cells[i][j].hasBottomWall = false;
      this.cells[iNext][jNext].hasTopWall = false;
    } else if (iNext < i) {
      this.cells[i][j].hasTopWall = false;
      this.cells[iNext][jNext].hasBottomWall = false;
    } else if (jNext > j) {
      this.cells[i][j].hasRightWall = false;
      this.cells[iNext][jNext].hasLeftWall = false;
    } else if (jNext < j) {
      this.cells[i][j].hasLeftWall = false;
      this.cells[iNext][jNext].hasRightWall = false;
    }

    this.drawCell(i, j);
    await this.animate(30);
    this.drawCell(iNext, jNext);
    await this.animate(30);
  }

  async breakWallsRecursive(i, j) {
    this.cells[i][j].visited = true;

    while (this.process) {
      let neighbours = [];
      if (i > 0 && !this.cells[i - 1][j].visited) {
        neighbours.push([i - 1, j]);
      }
      if (i < this.numRows - 1 && !this.cells[i + 1][j].visited) {
        neighbours.push([i + 1, j]);
      }
      if (j > 0 && !this.cells[i][j - 1].visited) {
        neighbours.push([i, j - 1]);
      }
      if (j < this.numCols - 1 && !this.cells[i][j + 1].visited) {
        neighbours.push([i, j + 1]);
      }

      if (neighbours.length === 0) {
        break;
      }

      const [iNext, jNext] = neighbours[Math.floor(Math.random() * neighbours.length)];
      await this.breakWalls(i, j, iNext, jNext);
      await this.breakWallsRecursive(iNext, jNext);
    }
  }

  resetCellsVisited() {
    for (let row = 0; row < this.numRows; row++) {
      for (let col = 0; col < this.numCols; col++) {
        if (!this.process) {
          return;
        }
        this.cells[row][col].visited = false;
      }
    }
  }

  solve() {
    this.resetCellsVisited();
    this.restoreCanvas();
    return this.solveRecursive(this.start[0], this.start[1]);
  }

  async solveRecursive(i, j) {
    this.cells[i][j].visited = true;
    await this.animate(100);

    if ((i === this.end[0] && j === this.end[1]) || !this.process) {
      return true;
    }

    let deadEnd = 0;

    // Down
    if (i < this.numRows - 1 && !this.cells[i + 1][j].visited && !this.cells[i][j].hasBottomWall) {
      this.cells[i][j].drawMove(this.cells[i + 1][j]);
      if (await this.solveRecursive(i + 1, j)) {
        return true;
      } else {
        this.cells[i][j].drawMove(this.cells[i + 1][j], true);
      }
    } else {
      deadEnd++;
    }

    // Right
    if (j < this.numCols - 1 && !this.cells[i][j + 1].visited && !this.cells[i][j].hasRightWall) {
      this.cells[i][j].drawMove(this.cells[i][j + 1]);
      if (await this.solveRecursive(i, j + 1)) {
        return true;
      } else {
        this.cells[i][j].drawMove(this.cells[i][j + 1], true);
      }
    } else {
      deadEnd++;
    }

    // Left
    if (j > 0 && !this.cells[i][j - 1].visited && !this.cells[i][j - 1].hasRightWall) {
      this.cells[i][j].drawMove(this.cells[i][j - 1]);
      if (await this.solveRecursive(i, j - 1)) {
        return true;
      } else {
        this.cells[i][j].drawMove(this.cells[i][j - 1], true);
      }
    } else {
      deadEnd++;
    }

    // Up
    if (i > 0 && !this.cells[i - 1][j].visited && !this.cells[i - 1][j].hasBottomWall) {
      this.cells[i][j].drawMove(this.cells[i - 1][j]);
      if (await this.solveRecursive(i - 1, j)) {
        return true;
      } else {
        this.cells[i][j].drawMove(this.cells[i - 1][j], true);
      }
    } else {
      deadEnd++;
    }

    if (deadEnd === 4 && this.process) {
      this.cells[i][j].drawEmoji("💥");
    }

    // Backtrack as no path found
    return false;
  }
  


}


