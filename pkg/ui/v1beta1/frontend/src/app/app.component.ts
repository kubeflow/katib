import { Component } from '@angular/core';
import { environment } from '@app/environment';
import { MatIconRegistry } from '@angular/material/icon';
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
})
export class AppComponent {
  constructor(
    private matIconRegistry: MatIconRegistry,
    private domSanitizer: DomSanitizer,
  ) {
    this.matIconRegistry.addSvgIcon(
      `pipeline-centered`,
      domSanitizer.bypassSecurityTrustResourceUrl(
        `${
          this.env.production
            ? '/katib/static/assets/pipeline-centered.svg'
            : '/assets/pipeline-centered.svg'
        }`,
      ),
    );
  }

  title = 'frontend';

  env = environment;
}
