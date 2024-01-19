class Point {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }
}

class Line {
  constructor(start, end) {
    this.start = start;
    this.end = end;
  }

  draw(ctx, fillColor = 'black', lineWidth = 2) {
    ctx.beginPath();
    ctx.moveTo(this.start.x, this.start.y);
    ctx.lineTo(this.end.x, this.end.y);
    ctx.strokeStyle = fillColor;
    ctx.lineWidth = lineWidth;
    ctx.stroke();
  }
}

class Cell {
  constructor(topLeft, bottomRight, canvas) {
    this._x1 = topLeft.x;
    this._y1 = topLeft.y;
    this._x2 = bottomRight.x;
    this._y2 = bottomRight.y;
    this.hasLeftWall = true;
    this.hasRightWall = true;
    this.hasTopWall = true;
    this.hasBottomWall = true;
    this.canvas = canvas;
    this.ctx = this.canvas.getContext('2d');
    this.visited = false;
  }

  draw() {
    const bgColor = "#ccc"; //window.getComputedStyle(this.canvas).backgroundColor;

    const drawLine = (line, color = 'black') => line.draw(this.ctx, color);
    const eraseLine = (line, color = bgColor) => line.draw(this.ctx, color, 3);

    const leftLine = new Line(new Point(this._x1, this._y1), new Point(this._x1, this._y2));
    this.hasLeftWall ? drawLine(leftLine) : eraseLine(leftLine, bgColor);

    const rightLine = new Line(new Point(this._x2, this._y1), new Point(this._x2, this._y2));
    this.hasRightWall ? drawLine(rightLine) : eraseLine(rightLine, bgColor);

    const topLine = new Line(new Point(this._x1, this._y1), new Point(this._x2, this._y1));
    this.hasTopWall ? drawLine(topLine) : eraseLine(topLine, bgColor);

    const bottomLine = new Line(new Point(this._x1, this._y2), new Point(this._x2, this._y2));
    this.hasBottomWall ? drawLine(bottomLine) : eraseLine(bottomLine, bgColor);
  }

  drawMove(toCell, undo = false) {
    const midPoint = (a, b) => (a + b) / 2;
    const point1 = new Point(midPoint(this._x1, this._x2), midPoint(this._y1, this._y2));
    const point2 = new Point(midPoint(toCell._x1, toCell._x2), midPoint(toCell._y1, toCell._y2));
    const line = new Line(point1, point2);
    if (undo) {
      line.draw(this.ctx, "#aaa");
      setTimeout(() => {}, 100); // Simulate a delay
    } else {
      line.draw(this.ctx, "red", 3);
    }
  }

  min(a, b) {
    return a < b ? a : b;
  }

  drawEmoji(emoji, fontSize = 40) {
    const x = (this._x1 + this._x2) / 2;
    const y = (this._y1 + this._y2) / 2;

    this.ctx.font = `${this.min(this.min(fontSize, Math.abs(this._x2 - this._x1 -6)), Math.abs(this._y2 - this._y1 -6))}px "Noto Color Emoji", "Apple Color Emoji", "Segoe UI Emoji", sans-serif`;
    this.ctx.textAlign = "center"; 
    this.ctx.textBaseline = "middle"; 

    this.ctx.fillText(emoji, x, y);
  }
  
}

// Example usage
// document.addEventListener('DOMContentLoaded', () => {
//   const canvas = document.getElementById('myCanvas');
  
//   if (!canvas) {
//     console.error("Canvas element not found");
//   }
// });
