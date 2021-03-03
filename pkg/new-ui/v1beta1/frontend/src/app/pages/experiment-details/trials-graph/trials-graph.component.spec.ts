import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialsGraphComponent } from './trials-graph.component';

describe('TrialsGraphComponent', () => {
  let component: TrialsGraphComponent;
  let fixture: ComponentFixture<TrialsGraphComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialsGraphComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialsGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
