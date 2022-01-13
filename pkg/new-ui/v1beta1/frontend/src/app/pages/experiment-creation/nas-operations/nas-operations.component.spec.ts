import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { NasOperationsComponent } from './nas-operations.component';

describe('NasOperationsComponent', () => {
  let component: NasOperationsComponent;
  let fixture: ComponentFixture<NasOperationsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [NasOperationsComponent],
      }).compileComponents();
    }),
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(NasOperationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
