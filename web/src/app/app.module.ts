import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ServiceWorkerModule } from '@angular/service-worker';
import { UxwbServicesModule } from '@uxwb/services';
import { ToastrModule } from 'ngx-toastr';
import { GroupsService } from 'src/app/pages/groups/groups.service';
import { UsersService } from 'src/app/pages/users/users.service';
import { AlertService } from 'src/app/services/alert.service';
import { ErrorInterceptor } from 'src/app/services/error-interceptor.service';
import { environment } from '../environments/environment';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { Error404Component } from './pages/error404/error404.component';
import { FilesComponent } from './pages/files/files.component';
import { GroupsComponent } from './pages/groups/groups.component';
import { HomeComponent } from './pages/home/home.component';
import { StoragesComponent } from './pages/storages/storages.component';
import { UsersComponent } from './pages/users/users.component';

@NgModule({
  declarations: [
    AppComponent,
    UsersComponent,
    GroupsComponent,
    StoragesComponent,
    FilesComponent,
    Error404Component,
    HomeComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    FormsModule,
    UxwbServicesModule.forRoot({ ajaxPrefixUrl:'/api/v1' }),
    ServiceWorkerModule.register('ngsw-worker.js', {
      enabled: environment.production,
      registrationStrategy: 'registerWhenStable:30000',
    }),
    ToastrModule.forRoot({
      preventDuplicates: true,
      progressBar: true,
      positionClass: 'toast-bottom-right',
    }),
  ],
  providers: [
    AlertService,
    UsersService,
    GroupsService,
    { provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor, multi: true },
  ],
  bootstrap: [AppComponent],
})
export class AppModule { }
