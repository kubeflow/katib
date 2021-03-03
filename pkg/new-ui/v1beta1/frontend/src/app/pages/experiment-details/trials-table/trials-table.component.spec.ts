import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialsTableComponent } from './trials-table.component';

describe('TrialsTableComponent', () => {
  let component: TrialsTableComponent;
  let fixture: ComponentFixture<TrialsTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialsTableComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
