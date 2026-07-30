'use client';

import React, { useMemo, useState } from 'react';
import styled, { useTheme } from 'styled-components';
import { useMutation } from '@apollo/client/react';
import { useConfig, useGoEnterpriseOffsets } from '@/hooks';
import type { GoEnterpriseOffsetsUpdateCheck } from '@/hooks';
import { CHECK_GO_ENTERPRISE_OFFSETS_UPDATES, UPDATE_GO_ENTERPRISE_OFFSETS } from '@/graphql';
import { GO_OFFSETS_PUBLIC_URL } from '@/utils';
import { ChevronDownIcon, ChevronRightIcon, GoLogo, RefreshIcon } from '@odigos/ui-kit/icons';
import { useNotificationStore } from '@odigos/ui-kit/store';
import { StatusType } from '@odigos/ui-kit/types';
import {
  Button,
  ButtonVariants,
  FlexColumn,
  FlexRow,
  Loader,
  PageContent,
  Search,
  Tag,
  TagVariants,
  Typography,
  TypographyColor,
  TypographySize,
  TypographyVariants,
} from '@odigos/ui-kit/components';

const Header = styled(FlexColumn)`
  margin-bottom: 8px;
`;

const Meta = styled(Typography)`
  opacity: 0.8;
`;

const AboutOffsets = styled(FlexColumn)`
  max-width: 720px;
  padding: 12px 14px;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  background: ${({ theme }) => theme.colors.dropdown_bg};
`;

const AboutOffsetsList = styled.ul`
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 4px;

  li {
    color: ${({ theme }) => theme.text.grey};
    font-size: ${({ theme }) => theme.v2.text.size.xs}px;
    line-height: 1.45;
  }
`;

const Toolbar = styled(FlexRow)`
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
`;

const List = styled(FlexColumn)`
  flex: 1;
  min-height: 0;
  overflow: auto;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
`;

const CheckBanner = styled(FlexRow)<{ $neutral?: boolean }>`
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  padding: 12px 16px;
  border: 1px solid ${({ theme, $neutral }) => ($neutral ? theme.colors.border : `${theme.text.success}55`)};
  border-left: 4px solid ${({ theme, $neutral }) => ($neutral ? theme.colors.border : theme.text.success)};
  border-radius: 8px;
  background: ${({ theme, $neutral }) => ($neutral ? theme.colors.dropdown_bg : `${theme.text.success}18`)};
`;

const CheckBannerTitle = styled(Typography)<{ $neutral?: boolean }>`
  color: ${({ theme, $neutral }) => ($neutral ? theme.colors.secondary : theme.text.success)};
  font-weight: 600;
`;

const ModuleHeader = styled.button<{ $expanded: boolean; $hasUpdates?: boolean }>`
  all: unset;
  box-sizing: border-box;
  display: grid;
  grid-template-columns: 24px minmax(0, 2fr) minmax(120px, 0.7fr) minmax(120px, 0.7fr);
  gap: 12px;
  align-items: center;
  width: 100%;
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme, $expanded, $hasUpdates }) =>
    $hasUpdates ? `${theme.colors.dark_green}22` : $expanded ? theme.colors.dropdown_bg_2 : 'transparent'};

  &:hover {
    background: ${({ theme, $hasUpdates }) => ($hasUpdates ? `${theme.colors.dark_green}33` : theme.colors.dropdown_bg_2)};
  }
`;

const ListHeader = styled.div`
  display: grid;
  grid-template-columns: 24px minmax(0, 2fr) minmax(120px, 0.7fr) minmax(120px, 0.7fr);
  gap: 12px;
  align-items: center;
  width: 100%;
  padding: 10px 16px;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  position: sticky;
  top: 0;
  background: ${({ theme }) => theme.colors.dropdown_bg};
  z-index: 1;
`;

const ExpandedPanel = styled(FlexColumn)<{ $hasUpdates?: boolean }>`
  padding: 8px 16px 16px 52px;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme, $hasUpdates }) => ($hasUpdates ? `${theme.colors.dark_green}0F` : theme.colors.dropdown_bg)};
`;

const MinorVersionRow = styled.div<{ $isNew?: boolean }>`
  display: grid;
  grid-template-columns: 72px minmax(0, 1fr);
  gap: 12px;
  align-items: start;
  padding: 8px 0;
  ${({ $isNew, theme }) => ($isNew ? `background: ${theme.colors.dark_green}14; border-radius: 6px; padding: 8px;` : '')}
`;

const VersionsWrap = styled(FlexRow)`
  flex-wrap: wrap;
`;

const EmptyState = styled(FlexColumn)`
  padding: 48px 16px;
  align-items: center;
`;

const Truncate = styled(Typography)`
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
`;

type DisplayVersion = {
  version: string;
  isNew?: boolean;
};

type DisplayMinor = {
  minorVersion: string;
  isNew?: boolean;
  versions: DisplayVersion[];
};

type DisplayModule = {
  module: string;
  isNew?: boolean;
  hasUpdates?: boolean;
  minVersion: string;
  maxVersion: string;
  minorVersions: DisplayMinor[];
};

type CheckUpdatesResponse = {
  checkGoEnterpriseOffsetsUpdates: GoEnterpriseOffsetsUpdateCheck;
};

const moduleHasUpdates = (mod: DisplayModule) =>
  !!mod.isNew || !!mod.hasUpdates || mod.minorVersions.some((minor) => !!minor.isNew || minor.versions.some((v) => !!v.isNew));

export default function Page() {
  const theme = useTheme();
  const { isReadonly } = useConfig();
  const { goEnterpriseOffsets, loading, error, refetch } = useGoEnterpriseOffsets();
  const { addNotification } = useNotificationStore();
  const [search, setSearch] = useState('');
  const [expandedModules, setExpandedModules] = useState<Set<string>>(new Set());
  const [updateCheck, setUpdateCheck] = useState<GoEnterpriseOffsetsUpdateCheck | null>(null);
  const [checkCompleted, setCheckCompleted] = useState(false);

  const [updateOffsets, { loading: mutating }] = useMutation(UPDATE_GO_ENTERPRISE_OFFSETS);
  const [checkOffsets, { loading: checkingMutation }] = useMutation<CheckUpdatesResponse>(CHECK_GO_ENTERPRISE_OFFSETS_UPDATES);
  const [fetching, setFetching] = useState(false);
  const [checking, setChecking] = useState(false);
  const updating = fetching || mutating;
  const checkingUpdates = checking || checkingMutation;

  const displayMods: DisplayModule[] = useMemo(() => {
    if (updateCheck?.hasUpdates) {
      return updateCheck.mods.map((mod) => {
        const minorVersions = mod.minorVersions.map((minor) => ({
          minorVersion: minor.minorVersion,
          isNew: minor.isNew,
          versions: minor.versions.map((v) => ({ version: v.version, isNew: v.isNew })),
        }));
        const hasUpdates =
          mod.isNew || minorVersions.some((minor) => minor.isNew || minor.versions.some((v) => v.isNew));
        return {
          module: mod.module,
          isNew: mod.isNew,
          hasUpdates,
          minVersion: mod.minVersion,
          maxVersion: mod.maxVersion,
          minorVersions,
        };
      });
    }

    return (goEnterpriseOffsets?.mods ?? []).map((mod) => ({
      module: mod.module,
      minVersion: mod.minVersion,
      maxVersion: mod.maxVersion,
      minorVersions: mod.minorVersions.map((minor) => ({
        minorVersion: minor.minorVersion,
        versions: minor.versions.map((version) => ({ version })),
      })),
    }));
  }, [goEnterpriseOffsets?.mods, updateCheck]);

  const filteredMods = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return displayMods;
    return displayMods.filter(
      (mod) =>
        mod.module.toLowerCase().includes(q) ||
        mod.minVersion.toLowerCase().includes(q) ||
        mod.maxVersion.toLowerCase().includes(q) ||
        mod.minorVersions.some(
          (minor) =>
            minor.minorVersion.toLowerCase().includes(q) ||
            minor.versions.some((v) => v.version.toLowerCase().includes(q)),
        ),
    );
  }, [displayMods, search]);

  const timestampLabel = goEnterpriseOffsets?.timestamp
    ? new Date(goEnterpriseOffsets.timestamp).toLocaleString()
    : undefined;

  const toggleModule = (moduleName: string) => {
    setExpandedModules((current) => {
      const next = new Set(current);
      if (next.has(moduleName)) {
        next.delete(moduleName);
      } else {
        next.add(moduleName);
      }
      return next;
    });
  };

  const fetchPublicOffsets = async () => {
    const response = await fetch(GO_OFFSETS_PUBLIC_URL);
    if (!response.ok) {
      throw new Error(`Failed to fetch offsets (${response.status})`);
    }
    return response.text();
  };

  const onCheckForUpdates = async () => {
    setChecking(true);
    try {
      const content = await fetchPublicOffsets();
      const result = await checkOffsets({ variables: { content } });
      const check = result.data?.checkGoEnterpriseOffsetsUpdates;
      if (!check) {
        throw new Error('No response from checkGoEnterpriseOffsetsUpdates');
      }
      setUpdateCheck(check);
      setCheckCompleted(true);
      if (check.hasUpdates) {
        setExpandedModules(
          new Set(
            check.mods
              .filter(
                (mod) =>
                  mod.isNew ||
                  mod.minorVersions.some((minor) => minor.isNew || minor.versions.some((v) => v.isNew)),
              )
              .map((mod) => mod.module),
          ),
        );
      }
    } catch (err) {
      setUpdateCheck(null);
      setCheckCompleted(false);
      addNotification({
        type: StatusType.Error,
        message: err instanceof Error ? err.message : 'Failed to check for go offset updates.',
      });
    } finally {
      setChecking(false);
    }
  };

  const onUpdateOffsets = async () => {
    setFetching(true);
    try {
      const content = await fetchPublicOffsets();
      await updateOffsets({ variables: { content } });
      setUpdateCheck(null);
      setCheckCompleted(false);
      await refetch();
      addNotification({ type: StatusType.Success, message: 'Go offsets updated from the latest public manifest.' });
    } catch (err) {
      addNotification({
        type: StatusType.Error,
        message: err instanceof Error ? err.message : 'Failed to update go offsets.',
      });
    } finally {
      setFetching(false);
    }
  };

  const dismissCheckResult = () => {
    setUpdateCheck(null);
    setCheckCompleted(false);
  };

  const checkBannerTitle = !updateCheck?.hasUpdates ? 'Go offsets are up to date' : 'New update to offsets available';
  const checkBannerDetail = !updateCheck?.hasUpdates
    ? 'No new versions found in the public manifest.'
    : 'Newer library versions are available. They are highlighted in the table below.';

  return (
    <PageContent $gap={16}>
      <Header $gap={8}>
        <FlexRow $gap={12} $alignItems='center' $justifyContent='space-between' $wrap='wrap'>
          <FlexRow $gap={12} $alignItems='center'>
            <GoLogo />
            <Typography variant={TypographyVariants.H2} size={TypographySize.XL}>
              Go Enterprise Offsets
            </Typography>
          </FlexRow>
          <FlexRow $gap={8} $alignItems='center' $wrap='wrap'>
            <Button
              data-id='check-go-offsets-updates'
              label='Check for updates'
              variant={ButtonVariants.Secondary}
              loading={checkingUpdates}
              disabled={checkingUpdates || updating}
              onClick={onCheckForUpdates}
            />
          </FlexRow>
        </FlexRow>
        <Meta variant={TypographyVariants.P} size={TypographySize.S} color={TypographyColor.Secondary}>
          {error
            ? error.message
            : timestampLabel
              ? `Offsets generated ${timestampLabel} · ${filteredMods.length} modules`
              : loading
                ? 'Loading go offsets…'
                : `${filteredMods.length} modules`}
        </Meta>
        <AboutOffsets $gap={6}>
          <Typography variant={TypographyVariants.P} size={TypographySize.S}>
            What are Go offsets?
          </Typography>
          <AboutOffsetsList>
            <li>Libraries and versions supported for automatic Go instrumentation.</li>
            <li>Other modules can still be instrumented with a custom instrumentation rule.</li>
            <li>If a version is missing here, it will not produce traces.</li>
          </AboutOffsetsList>
        </AboutOffsets>
      </Header>

      {checkCompleted && updateCheck && (
        <CheckBanner $neutral={!updateCheck.hasUpdates}>
          <FlexColumn $gap={2}>
            <CheckBannerTitle $neutral={!updateCheck.hasUpdates} variant={TypographyVariants.P} size={TypographySize.S}>
              {checkBannerTitle}
            </CheckBannerTitle>
            <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
              {checkBannerDetail}
            </Typography>
          </FlexColumn>
          <FlexRow $gap={8} $alignItems='center' $wrap='wrap'>
            {updateCheck.hasUpdates && (
              <Button
                data-id='update-go-offsets'
                label='Update offsets'
                leftIcon={RefreshIcon}
                variant={ButtonVariants.Primary}
                loading={updating}
                disabled={isReadonly || updating || checkingUpdates}
                onClick={onUpdateOffsets}
              />
            )}
            <Button data-id='dismiss-go-offsets-check' label='Dismiss' variant={ButtonVariants.Text} onClick={dismissCheckResult} />
          </FlexRow>
        </CheckBanner>
      )}

      <Toolbar $gap={12}>
        <Search data-id='go-offsets-search' value={search} onChange={setSearch} placeholder='Search modules or versions' width={360} />
      </Toolbar>

      <List>
        <ListHeader>
          <span />
          <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
            Module
          </Typography>
          <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
            Min version
          </Typography>
          <Typography variant={TypographyVariants.Span} size={TypographySize.XS} color={TypographyColor.Secondary}>
            Max version
          </Typography>
        </ListHeader>

        {loading && (
          <EmptyState $gap={12}>
            <Loader withSpinnerOld scaleSpinnerOld={1.5} />
          </EmptyState>
        )}

        {!loading && filteredMods.length === 0 && (
          <EmptyState $gap={8}>
            <Typography variant={TypographyVariants.P} size={TypographySize.M}>
              No go offsets
            </Typography>
            <Typography variant={TypographyVariants.P} size={TypographySize.S} color={TypographyColor.Secondary}>
              {error ? error.message : 'No modules found in the odigos-go-offsets ConfigMap.'}
            </Typography>
          </EmptyState>
        )}

        {!loading &&
          filteredMods.map((mod) => {
            const expanded = expandedModules.has(mod.module);
            const hasUpdates = moduleHasUpdates(mod);
            return (
              <React.Fragment key={mod.module}>
                <ModuleHeader
                  $expanded={expanded}
                  $hasUpdates={hasUpdates}
                  onClick={() => toggleModule(mod.module)}
                  type='button'
                >
                  {expanded ? <ChevronDownIcon fill={theme.colors.secondary} /> : <ChevronRightIcon fill={theme.colors.secondary} />}
                  <FlexRow $gap={8} $alignItems='center'>
                    <div title={mod.module} style={{ minWidth: 0, flex: 1 }}>
                      <Truncate variant={TypographyVariants.Span} size={TypographySize.S}>
                        {mod.module}
                      </Truncate>
                    </div>
                    {hasUpdates && <Tag variant={TagVariants.Green} label='New' />}
                  </FlexRow>
                  <Typography variant={TypographyVariants.Span} size={TypographySize.S}>
                    {mod.minVersion || '—'}
                  </Typography>
                  <Typography variant={TypographyVariants.Span} size={TypographySize.S}>
                    {mod.maxVersion || '—'}
                  </Typography>
                </ModuleHeader>
                {expanded && (
                  <ExpandedPanel $gap={0} $hasUpdates={hasUpdates}>
                    {mod.minorVersions.length === 0 ? (
                      <Typography variant={TypographyVariants.Span} size={TypographySize.S} color={TypographyColor.Secondary}>
                        No versions
                      </Typography>
                    ) : (
                      mod.minorVersions.map((minor) => (
                        <MinorVersionRow key={minor.minorVersion} $isNew={minor.isNew}>
                          <Typography variant={TypographyVariants.Span} size={TypographySize.S} color={TypographyColor.Secondary}>
                            {minor.minorVersion}
                          </Typography>
                          <VersionsWrap $gap={8}>
                            {minor.versions.map((version) => (
                              <Tag
                                key={version.version}
                                label={version.isNew ? `${version.version} · New` : version.version}
                                variant={version.isNew ? TagVariants.Green : TagVariants.Default}
                              />
                            ))}
                          </VersionsWrap>
                        </MinorVersionRow>
                      ))
                    )}
                  </ExpandedPanel>
                )}
              </React.Fragment>
            );
          })}
      </List>
    </PageContent>
  );
}
