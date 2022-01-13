import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { YamlModalComponent } from './yaml-modal.component';

describe('YamlModalComponent', () => {
  let component: YamlModalComponent;
  let fixture: ComponentFixture<YamlModalComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [YamlModalComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(YamlModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
