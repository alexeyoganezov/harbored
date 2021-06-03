import axios from 'axios';
import { OsfPresenter } from 'oldskull';
import { Presentation, PresentationModel } from './PresentationModel'
import { PresentationView } from './PresentationView'

export class PresentationPresenter extends OsfPresenter<PresentationModel, PresentationView> {
  model = new PresentationModel({
    id: 0,
    filename: '',
    isOnline: false,
    currentPageNumber: 0,
  });
  view = new PresentationView(this.model);
  ws?: WebSocket;
  presentations: Presentation[] = [];
  async beforeInit() {
    const response = await axios.get(`/api/presentations`);
    this.presentations = response.data;
    const onlinePresentation = this.presentations.find(el => el.isOnline);
    if (onlinePresentation) {
      this.model.set(onlinePresentation);
    }
    this.wsConnect();
  }
  afterInit() {
    // Cache fix for iOS Safari
    window.onpageshow = function(event: PageTransitionEvent) {
      if (event.persisted) {
        window.location.reload()
      }
    };
    // Reload on re-focus
    let hidden: string;
    let visibilityChange = '';
    if (typeof document.hidden !== 'undefined') {
      hidden = 'hidden';
      visibilityChange = 'visibilitychange';
    } else if (typeof (document as any).msHidden !== 'undefined') {
      hidden = 'msHidden';
      visibilityChange = 'msvisibilitychange';
    } else if (typeof (document as any).webkitHidden !== 'undefined') {
      hidden = 'webkitHidden';
      visibilityChange = 'webkitvisibilitychange';
    }
    document.addEventListener(visibilityChange, () => {
      // @ts-ignore
      if (this.ws && document[hidden]) {
        this.ws.close();
      } else {
        this.view.document.destroy();
        document.location.reload();
      }
    }, false);
  }
  wsConnect() {
    this.ws = new WebSocket(`ws://${location.host}/api/ws`);
    this.ws.onopen = () => {
      if (this.ws) {
        const message = {
          command: 'presentation:subscribe',
        };
        this.ws.send(JSON.stringify(message));
      }
    };
    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.command === 'presentation:start') {
        const presentation = <Presentation>this.presentations.find(el => el.id === message.payload.presentationId);
        this.model.set({ ...presentation, isOnline: true, currentPageNumber: 0 });
      } else if (message.command === 'presentation:slide-change') {
        this.model.set({ ...this.model.attrs, currentPageNumber: message.payload.pageNumber });
      } else if (message.command === 'presentation:stop') {
        if (message.presentationId === this.model.attrs.id) {
          this.model.set({ ...this.model.attrs, isOnline: false });
        }
      } else if (message.command === 'presentation:end') {
        this.view.handleEnd()
      }
    };
    this.ws.onerror = function(error) {
      console.error(error);
    };
    this.ws.onclose = () => {
      this.view.handleDisconnect();
    };
  }
}
