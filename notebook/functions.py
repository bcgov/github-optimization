import os
import pandas as pd

# Ignoring these repolint columns
cols_to_drop = [
    'javascript-package-metadata-exists', 
    'ruby-package-metadata-exists',
    'java-package-metadata-exists', 
    'python-package-metadata-exists',
    'objective-c-package-metadata-exists', 
    'swift-package-metadata-exists',
    'erlang-package-metadata-exists', 
    'elixir-package-metadata-exists',
    'license-detectable-by-licensee',
    'notice-file-exists'
]

def parse_file(file, org, base_df=pd.DataFrame()):
    directory = os.getcwd()
    dir_path = '{directory}/dat/{org}/{file}'.format(directory=directory, org=org, file=file)
    data = pd.read_csv(dir_path)
    if (base_df.empty):
      return data
    return data[data["Repository"].isin(base_df["Repository"])]

def left_merge(main, to_merge):
    return main.merge(to_merge, how="left", on="Repository", suffixes=('', '_DROP')).filter(regex='^(?!.*_DROP)')

def count_repositories(df, col_name):
  return df.pivot_table(index=['Repository'], aggfunc='size').to_frame().reset_index().rename(columns={0:col_name})

def format_repolint_results(df):
  formatted_df = df
  formatted_df.drop(cols_to_drop, axis=1, inplace=True)
  formatted_df.replace(inplace=True, to_replace="PASSED", value=True)
  formatted_df.replace(inplace=True, to_replace="NOT_PASSED*", value=False, regex=True)
  return formatted_df

def print_to_csv(df, org, filename):
  directory = os.getcwd()
  df.to_csv('{directory}/dat/{org}/{filename}'.format(directory=directory, org=org, filename=filename))