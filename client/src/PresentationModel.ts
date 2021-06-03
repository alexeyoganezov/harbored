import { OsfModel } from 'oldskull';

export type Presentation = {
  id: number;
  filename: string;
  isOnline: boolean;
  currentPageNumber: number;
}

export class PresentationModel extends OsfModel<Presentation> {}
