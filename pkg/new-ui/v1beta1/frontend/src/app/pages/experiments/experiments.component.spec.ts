import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ExperimentsComponent } from './experiments.component';

describe('ExperimentsComponent', () => {
  let component: ExperimentsComponent;
  let fixture: ComponentFixture<ExperimentsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExperimentsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExperimentsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
