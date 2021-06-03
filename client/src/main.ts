import { OsfApplication } from 'oldskull';
import { PresentationPresenter } from './PresentationPresenter';

export class HarboredApp extends OsfApplication {
  async init() {
    await this.mainRegion.show(new PresentationPresenter());
  }
}

const app = new HarboredApp('#app');
app.init();
