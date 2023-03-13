import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { Error404Component } from 'src/app/pages/error404/error404.component';
import { FilesComponent } from 'src/app/pages/files/files.component';
import { GroupsComponent } from 'src/app/pages/groups/groups.component';
import { HomeComponent } from 'src/app/pages/home/home.component';
import { StoragesComponent } from 'src/app/pages/storages/storages.component';
import { UsersComponent } from 'src/app/pages/users/users.component';

const routes: Routes = [
  { path:'home', component: HomeComponent },
  { path:'users', component: UsersComponent },
  { path:'groups', component: GroupsComponent },
  { path:'storages', component: StoragesComponent },
  { path:'files', component: FilesComponent },
  { path:'', redirectTo:'/home', pathMatch: 'full' },
  { path:'**', component: Error404Component },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule],
})
export class AppRoutingModule { }
