import { OsfModelView, OsfReference } from 'oldskull';
import { PresentationModel } from './PresentationModel';

const pdfjs = require('pdfjs-dist/es5/build/pdf');
require('pdfjs-dist/es5/build/pdf.worker.entry');
pdfjs.showPreviousViewOnLoad = false;

export class PresentationView extends OsfModelView<PresentationModel> {
  // PDF.js data has no type definitions
  document: any;
  page: any;
  renderTask: any;
  viewport: any;
  canvasContext: CanvasRenderingContext2D | null = null;
  devicePixelRatio: number = window.devicePixelRatio || 1;
  transform: any = [ devicePixelRatio, 0 , 0, devicePixelRatio, 0, 0];
  waitMessage = new OsfReference(this, '#wait');
  disconnectMessage = new OsfReference(this, '#disconnect');
  presentationEndMessage = new OsfReference(this, '#ended');
  canvasRef = new OsfReference<HTMLCanvasElement>(this, '#canvas');
  // Copy of model data to distinguish its changes
  // Should be removed on refactoring
  pageNumber: number;
  currentPresentationId: number;
  isOnline: boolean;
  constructor(model: PresentationModel) {
    super(model);
    this.pageNumber = model.attrs.currentPageNumber;
    this.currentPresentationId = model.attrs.id;
    this.isOnline = model.attrs.isOnline;
  }
  getHTML(): string {
    return `
      <div id="pdf-container">
        <div id="wait" class="popup">
          Waiting for start...
        </div>
        <div id="disconnect" class="popup">
          <p>
            Connecting...
          </p>
        </div>
        <div id="ended" class="popup">
          <p>
            The presentation is over
          </p>
        </div>
        <canvas id="canvas"></canvas>
      </div>
    `;
  }
  modelEvents = [
    {
      on: 'change',
      call: this.handleModelChange.bind(this),
    },
  ];
  protected async afterInit() {
    this.setMessagesHeight(window.innerHeight);
    if (this.model.attrs.isOnline) {
      this.showPresentation();
    } else {
      this.setWaiterVisibility(true);
    }
    // Re-render slide on phone reorientation
    window.addEventListener('orientationchange', async (event) => {
      const afterOrientationChange = () => {
        if (this.page) {
          this.renderSlide();
          window.scrollTo(0, 0);
        }
        window.removeEventListener('resize', afterOrientationChange);
      };
      window.addEventListener('resize', afterOrientationChange);
    });
  }
  private async showPresentation() {
    this.document = await pdfjs.getDocument({
      url: `/static/${this.model.attrs.filename}`,
    }).promise;
    this.page = await this.document.getPage(this.model.attrs.currentPageNumber + 1);
    this.renderSlide();
    this.setWaiterVisibility(false);
  }
  private renderSlide() {
    let scale = 1;
    let viewport = this.page.getViewport({ scale });
    const flag = window.innerWidth > window.innerHeight;
    const w = flag ? window.innerWidth : window.innerHeight;
    const h = flag ? window.innerHeight : window.innerWidth;
    if (viewport.width > w) {
      scale = w / this.page._pageInfo.view[2];
    }
    if (viewport.height > h) {
      const value = h / this.page._pageInfo.view[3];
      if (value < scale) {
        scale = value;
      }
    }
    viewport = this.page.getViewport({ scale });
    const canvas = this.canvasRef.get();
    if (!canvas) {
      throw new Error('Canvas element not found');
    }
    const canvasContext = canvas.getContext('2d');
    canvas.style.maxHeight = `${h}px`;
    this.setMessagesHeight(h);
    canvas.width = viewport.width * this.devicePixelRatio;
    canvas.height = viewport.height * this.devicePixelRatio;
    this.page.render({
      canvasContext,
      viewport,
      transform: this.transform,
    });
    this.canvasContext = canvasContext;
    this.viewport = viewport;
  }
  private setWaiterVisibility(status: boolean) {
    const el = <HTMLElement>this.waitMessage.get();
    if (status) {
      el.style.display = 'flex'
    } else {
      el.style.display = 'none'
    }
  }
  private async changePage(pageNumber: number) {
    if (this.renderTask) {
      await this.renderTask.promise;
    }
    this.page = await this.document.getPage(pageNumber + 1);
    this.renderTask = this.page.render({
      canvasContext: this.canvasContext,
      viewport: this.viewport,
      transform: this.transform,
    });
    this.renderTask.promise
      .then(() => {
        this.renderTask = null
      })
      .catch((err: Error) => {
        if (err.name !== 'RenderingCancelledException') {
          throw err;
        }
      })
  }
  private setMessagesHeight(height: number) {
    let value: string;
    if (window.innerWidth < window.innerHeight) {
      value = '100%';
    } else {
      value = `${height}px`;
    }
    const waitEl = <HTMLElement>this.waitMessage.get();
    waitEl.style.height = value;
    const discEl = <HTMLElement>this.disconnectMessage.get();
    discEl.style.height = value;
    const endedEl = <HTMLElement>this.presentationEndMessage.get();
    endedEl.style.height = value;
  }
  private handleModelChange(payload: unknown) {
    let m = <PresentationModel>payload;
    if (m.attrs.id !== this.currentPresentationId) {
      this.currentPresentationId = m.attrs.id;
      this.isOnline = m.attrs.isOnline;
      this.pageNumber = m.attrs.currentPageNumber;
      this.showPresentation();
    } else {
      if (m.attrs.currentPageNumber !== this.pageNumber) {
        this.pageNumber = m.attrs.currentPageNumber;
        this.changePage(this.pageNumber);
        this.setWaiterVisibility(false);
      }
      if (m.attrs.isOnline !== this.isOnline) {
        this.isOnline = this.model.attrs.isOnline;
        this.setWaiterVisibility(!this.model.attrs.isOnline)
      }
    }
  }
  public handleDisconnect() {
    const el = <HTMLElement>this.disconnectMessage.get();
    el.style.display = 'flex';
  }
  public handleEnd() {
    const el = <HTMLElement>this.presentationEndMessage.get();
    el.style.display = 'flex';
  }
}
