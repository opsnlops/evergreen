Name:           evergreen
Version:        %{_git_hash}
Release:        1%{?dist}
Summary:        A Distributed Continuous Integration System from MongoDB

License:        Apache
URL:            https://github.com/evergreen-ci/%{name}


Source0:         %{name}-%{_git_hash}.tar.gz
BuildRoot:      %{_tmppath}/%{name}-%{_git_hash}-build

BuildArch:      aarch64


BuildRequires:  make
BuildRequires:  go
BuildRequires:  git

# Golang apps are missing a few things that gcc adds. We don't need them.
%global _missing_build_ids_terminate_build 0
%global debug_package %{nil}

%description
A Distributed Continuous Integration System from MongoDB

%prep
%setup -n %{name}-%{_git_hash}

%build
/root/go/bin/bom generate --format json --name %{name} --output %{name}-%{_git_hash}.spdx.json .
make

%install
install -D -m 755 clients/%{_go_os}_%{_go_arch}/%{name} %{buildroot}%{_bindir}/%{name}
install -D -m 444 %{name}-%{_git_hash}.spdx.json %{buildroot}/var/lib/db/sbom/%{name}-%{_git_hash}.spdx.json


%files
%license LICENSE.md
%{_bindir}/%{name}
/var/lib/db/sbom/%{name}-%{_git_hash}.spdx.json

%changelog
* Wed May 10 2023 April White <april.white@mongodb.com> - 1.0-1
- First evergreen package
