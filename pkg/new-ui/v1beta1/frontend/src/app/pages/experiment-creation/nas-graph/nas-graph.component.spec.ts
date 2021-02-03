import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NasGraphComponent } from './nas-graph.component';

describe('NasGraphComponent', () => {
  let component: NasGraphComponent;
  let fixture: ComponentFixture<NasGraphComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [NasGraphComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NasGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
