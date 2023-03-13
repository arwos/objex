import { Component, OnInit } from '@angular/core';
import { GroupsService } from 'src/app/pages/groups/groups.service';
import { Group } from 'src/app/pages/groups/models';

@Component({
  selector: 'app-groups',
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent implements OnInit {

  list: Group[] = [];

  constructor(private readonly groupsService: GroupsService) { }

  ngOnInit(): void {
    this.groupsService.list().subscribe(value => this.list = value);
  }

}
