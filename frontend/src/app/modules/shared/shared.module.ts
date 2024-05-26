import {CommonModule, DatePipe} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import {CUSTOM_ELEMENTS_SCHEMA, NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {SharedRoutes} from './shared.routes';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {DateFormatPipe} from './pipes/date-format.pipe';
import {CustomMaterialModuleModule} from "../../custom-material-module/custom-material-module.module";
import {HomePageComponent} from "./pages/home-page/home-page.component";
import {SearchComponent} from './components/search/search.component';
import {NgxStarsModule} from "ngx-stars";
import {CarouselModule} from 'primeng/carousel';
import {AccommodationViewComponent} from './components/accomodation-view/accommodation-view.component';

@NgModule({
  declarations: [
    HomePageComponent,
    DateFormatPipe,
    SearchComponent,
    AccommodationViewComponent,
    AccommodationViewComponent,
    AccommodationViewComponent
  ],
  exports: [
    DateFormatPipe,
    AccommodationViewComponent
  ],
  providers: [
    DatePipe,
  ],
  imports: [
    CommonModule,
    CustomMaterialModuleModule,
    FormsModule,
    HttpClientModule,
    RouterModule.forChild(SharedRoutes),
    ReactiveFormsModule,
    NgxStarsModule,
    CarouselModule,
  ],
  schemas: [ CUSTOM_ELEMENTS_SCHEMA ]
})
export class SharedModule {}
