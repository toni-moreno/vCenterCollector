import { SnmpcollectorPage } from './app.po';

describe('vcentercollector App', function() {
  let page: vcentercollectorPage;

  beforeEach(() => {
    page = new vcentercollectorPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
