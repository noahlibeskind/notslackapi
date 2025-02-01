create table "public"."channels" (
    "id" uuid not null default gen_random_uuid(),
    "workspace_id" uuid not null,
    "name" text not null
);


alter table "public"."channels" enable row level security;

create table "public"."message" (
    "id" uuid not null default gen_random_uuid(),
    "sender_id" uuid not null,
    "posted_at" timestamp with time zone not null,
    "content" text not null,
    "channel_id" uuid not null
);


alter table "public"."message" enable row level security;

create table "public"."users" (
    "id" uuid not null default gen_random_uuid(),
    "first_name" text not null,
    "last_name" text not null
);


alter table "public"."users" enable row level security;

create table "public"."workspace_members" (
    "id" uuid not null default gen_random_uuid(),
    "member_id" uuid not null,
    "workspace_id" uuid not null
);


alter table "public"."workspace_members" enable row level security;

create table "public"."workspaces" (
    "id" uuid not null default gen_random_uuid(),
    "name" text not null,
    "owner_id" uuid
);


alter table "public"."workspaces" enable row level security;

CREATE UNIQUE INDEX channels_pkey ON public.channels USING btree (id);

CREATE UNIQUE INDEX message_pkey ON public.message USING btree (id);

CREATE UNIQUE INDEX us_pkey ON public.users USING btree (id);

CREATE UNIQUE INDEX workspace_members_pkey ON public.workspace_members USING btree (id);

CREATE UNIQUE INDEX workspaces_pkey ON public.workspaces USING btree (id);

alter table "public"."channels" add constraint "channels_pkey" PRIMARY KEY using index "channels_pkey";

alter table "public"."message" add constraint "message_pkey" PRIMARY KEY using index "message_pkey";

alter table "public"."users" add constraint "us_pkey" PRIMARY KEY using index "us_pkey";

alter table "public"."workspace_members" add constraint "workspace_members_pkey" PRIMARY KEY using index "workspace_members_pkey";

alter table "public"."workspaces" add constraint "workspaces_pkey" PRIMARY KEY using index "workspaces_pkey";

alter table "public"."channels" add constraint "channels_workspace_id_fkey" FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."channels" validate constraint "channels_workspace_id_fkey";

alter table "public"."message" add constraint "message_channel_id_fkey" FOREIGN KEY (channel_id) REFERENCES channels(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."message" validate constraint "message_channel_id_fkey";

alter table "public"."message" add constraint "message_sender_id_fkey" FOREIGN KEY (sender_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."message" validate constraint "message_sender_id_fkey";

alter table "public"."users" add constraint "users_id_fkey" FOREIGN KEY (id) REFERENCES auth.users(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."users" validate constraint "users_id_fkey";

alter table "public"."workspace_members" add constraint "workspace_members_member_id_fkey" FOREIGN KEY (member_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."workspace_members" validate constraint "workspace_members_member_id_fkey";

alter table "public"."workspace_members" add constraint "workspace_members_workspace_id_fkey" FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON UPDATE CASCADE ON DELETE CASCADE not valid;

alter table "public"."workspace_members" validate constraint "workspace_members_workspace_id_fkey";

alter table "public"."workspaces" add constraint "workspaces_owner_id_fkey" FOREIGN KEY (owner_id) REFERENCES users(id) not valid;

alter table "public"."workspaces" validate constraint "workspaces_owner_id_fkey";

grant delete on table "public"."channels" to "anon";

grant insert on table "public"."channels" to "anon";

grant references on table "public"."channels" to "anon";

grant select on table "public"."channels" to "anon";

grant trigger on table "public"."channels" to "anon";

grant truncate on table "public"."channels" to "anon";

grant update on table "public"."channels" to "anon";

grant delete on table "public"."channels" to "authenticated";

grant insert on table "public"."channels" to "authenticated";

grant references on table "public"."channels" to "authenticated";

grant select on table "public"."channels" to "authenticated";

grant trigger on table "public"."channels" to "authenticated";

grant truncate on table "public"."channels" to "authenticated";

grant update on table "public"."channels" to "authenticated";

grant delete on table "public"."channels" to "service_role";

grant insert on table "public"."channels" to "service_role";

grant references on table "public"."channels" to "service_role";

grant select on table "public"."channels" to "service_role";

grant trigger on table "public"."channels" to "service_role";

grant truncate on table "public"."channels" to "service_role";

grant update on table "public"."channels" to "service_role";

grant delete on table "public"."message" to "anon";

grant insert on table "public"."message" to "anon";

grant references on table "public"."message" to "anon";

grant select on table "public"."message" to "anon";

grant trigger on table "public"."message" to "anon";

grant truncate on table "public"."message" to "anon";

grant update on table "public"."message" to "anon";

grant delete on table "public"."message" to "authenticated";

grant insert on table "public"."message" to "authenticated";

grant references on table "public"."message" to "authenticated";

grant select on table "public"."message" to "authenticated";

grant trigger on table "public"."message" to "authenticated";

grant truncate on table "public"."message" to "authenticated";

grant update on table "public"."message" to "authenticated";

grant delete on table "public"."message" to "service_role";

grant insert on table "public"."message" to "service_role";

grant references on table "public"."message" to "service_role";

grant select on table "public"."message" to "service_role";

grant trigger on table "public"."message" to "service_role";

grant truncate on table "public"."message" to "service_role";

grant update on table "public"."message" to "service_role";

grant delete on table "public"."users" to "anon";

grant insert on table "public"."users" to "anon";

grant references on table "public"."users" to "anon";

grant select on table "public"."users" to "anon";

grant trigger on table "public"."users" to "anon";

grant truncate on table "public"."users" to "anon";

grant update on table "public"."users" to "anon";

grant delete on table "public"."users" to "authenticated";

grant insert on table "public"."users" to "authenticated";

grant references on table "public"."users" to "authenticated";

grant select on table "public"."users" to "authenticated";

grant trigger on table "public"."users" to "authenticated";

grant truncate on table "public"."users" to "authenticated";

grant update on table "public"."users" to "authenticated";

grant delete on table "public"."users" to "service_role";

grant insert on table "public"."users" to "service_role";

grant references on table "public"."users" to "service_role";

grant select on table "public"."users" to "service_role";

grant trigger on table "public"."users" to "service_role";

grant truncate on table "public"."users" to "service_role";

grant update on table "public"."users" to "service_role";

grant delete on table "public"."workspace_members" to "anon";

grant insert on table "public"."workspace_members" to "anon";

grant references on table "public"."workspace_members" to "anon";

grant select on table "public"."workspace_members" to "anon";

grant trigger on table "public"."workspace_members" to "anon";

grant truncate on table "public"."workspace_members" to "anon";

grant update on table "public"."workspace_members" to "anon";

grant delete on table "public"."workspace_members" to "authenticated";

grant insert on table "public"."workspace_members" to "authenticated";

grant references on table "public"."workspace_members" to "authenticated";

grant select on table "public"."workspace_members" to "authenticated";

grant trigger on table "public"."workspace_members" to "authenticated";

grant truncate on table "public"."workspace_members" to "authenticated";

grant update on table "public"."workspace_members" to "authenticated";

grant delete on table "public"."workspace_members" to "service_role";

grant insert on table "public"."workspace_members" to "service_role";

grant references on table "public"."workspace_members" to "service_role";

grant select on table "public"."workspace_members" to "service_role";

grant trigger on table "public"."workspace_members" to "service_role";

grant truncate on table "public"."workspace_members" to "service_role";

grant update on table "public"."workspace_members" to "service_role";

grant delete on table "public"."workspaces" to "anon";

grant insert on table "public"."workspaces" to "anon";

grant references on table "public"."workspaces" to "anon";

grant select on table "public"."workspaces" to "anon";

grant trigger on table "public"."workspaces" to "anon";

grant truncate on table "public"."workspaces" to "anon";

grant update on table "public"."workspaces" to "anon";

grant delete on table "public"."workspaces" to "authenticated";

grant insert on table "public"."workspaces" to "authenticated";

grant references on table "public"."workspaces" to "authenticated";

grant select on table "public"."workspaces" to "authenticated";

grant trigger on table "public"."workspaces" to "authenticated";

grant truncate on table "public"."workspaces" to "authenticated";

grant update on table "public"."workspaces" to "authenticated";

grant delete on table "public"."workspaces" to "service_role";

grant insert on table "public"."workspaces" to "service_role";

grant references on table "public"."workspaces" to "service_role";

grant select on table "public"."workspaces" to "service_role";

grant trigger on table "public"."workspaces" to "service_role";

grant truncate on table "public"."workspaces" to "service_role";

grant update on table "public"."workspaces" to "service_role";



