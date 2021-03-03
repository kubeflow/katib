import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ParamsListComponent } from './params-list.component';

describe('ParamsListComponent', () => {
  let component: ParamsListComponent;
  let fixture: ComponentFixture<ParamsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ParamsListComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ParamsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
